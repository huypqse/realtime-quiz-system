package logic

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"realtime_quiz_system/api"
	"realtime_quiz_system/internal/dao"
	"realtime_quiz_system/internal/model/do"
	"realtime_quiz_system/internal/model/entity"
	"realtime_quiz_system/internal/service"
	"realtime_quiz_system/utility"
	"time"

	"github.com/gogf/gf/v2/os/glog"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type SessionServiceImpl struct {
	logger            *glog.Logger
	redisClient       *redis.Client
	centrifugoService service.CentrifugoService
}

func NewSessionService(logger *glog.Logger, redisClient *redis.Client, centrifugoService service.CentrifugoService) service.SessionService {
	return &SessionServiceImpl{
		logger:            logger,
		redisClient:       redisClient,
		centrifugoService: centrifugoService,
	}
}

func (s *SessionServiceImpl) CreateSession(ctx context.Context, userId string, req *api.CreateSessionReq) (*api.CreateSessionRes, error) {
	s.logger.Info(ctx, "Creating session", "userId", userId, "quizId", req.QuizId)

	// Parse user ID
	hostId, err := uuid.Parse(userId)
	if err != nil {
		return nil, errors.New("invalid user ID")
	}

	// Parse quiz ID
	quizId, err := uuid.Parse(req.QuizId)
	if err != nil {
		return nil, errors.New("invalid quiz ID")
	}

	// Verify quiz exists
	var quiz entity.Quizzes
	err = dao.Quizzes.Ctx(ctx).Where("id", quizId).Scan(&quiz)
	if err != nil || quiz.Id.String() == "00000000-0000-0000-0000-000000000000" {
		return nil, errors.New("quiz not found")
	}

	// Set default max participants
	maxParticipants := req.MaxParticipants
	if maxParticipants == 0 {
		maxParticipants = 50
	}

	// Generate unique session code
	sessionCode, err := s.generateUniqueSessionCode(ctx)
	if err != nil {
		s.logger.Error(ctx, "Failed to generate session code", "error", err)
		return nil, errors.New("failed to generate session code")
	}

	// Create session
	sessionId, err := uuid.NewV7()
	if err != nil {
		return nil, errors.New("failed to generate session ID")
	}

	_, err = dao.QuizSessions.Ctx(ctx).Data(do.QuizSessions{
		Id:                   sessionId,
		QuizId:               quizId,
		SessionCode:          sessionCode,
		HostId:               hostId,
		Status:               "waiting",
		CurrentQuestionIndex: 0,
		CurrentQuestionId:    nil,
		MaxParticipants:      maxParticipants,
	}).Insert()
	if err != nil {
		s.logger.Error(ctx, "Failed to create session", "error", err)
		return nil, errors.New("failed to create session")
	}

	// Cache session info in Redis
	err = s.cacheSessionInfo(ctx, sessionId.String(), map[string]interface{}{
		"quiz_id":            quizId.String(),
		"status":             "waiting",
		"current_question":   "",
		"participants_count": 0,
		"started_at":         "",
	})
	if err != nil {
		s.logger.Warning(ctx, "Failed to cache session info", "error", err)
		// Continue even if Redis caching fails
	}

	// Get created session with quiz info
	var session entity.QuizSessions
	err = dao.QuizSessions.Ctx(ctx).Where("id", sessionId).Scan(&session)
	if err != nil {
		return nil, errors.New("failed to retrieve created session")
	}

	s.logger.Info(ctx, "Session created successfully", "sessionId", sessionId.String(), "sessionCode", sessionCode)

	return &api.CreateSessionRes{
		SessionId:       sessionId.String(),
		SessionCode:     sessionCode,
		QuizId:          quiz.Id.String(),
		QuizTitle:       quiz.Title,
		TotalQuestions:  quiz.TotalQuestions,
		Status:          "waiting",
		HostId:          hostId.String(),
		MaxParticipants: maxParticipants,
		CreatedAt:       session.CreatedAt.String(),
	}, nil
}

func (s *SessionServiceImpl) JoinSession(ctx context.Context, userId string, req *api.JoinSessionReq) (*api.JoinSessionRes, error) {
	s.logger.Info(ctx, "User joining session", "userId", userId, "sessionCode", req.SessionCode)

	// Parse user ID
	uid, err := uuid.Parse(userId)
	if err != nil {
		return nil, errors.New("invalid user ID")
	}

	// Normalize session code
	normalizedCode := utility.NormalizeSessionCode(req.SessionCode)

	// Find session by code
	var session entity.QuizSessions
	err = dao.QuizSessions.Ctx(ctx).Where("session_code", normalizedCode).Scan(&session)
	if err != nil || session.Id.String() == "00000000-0000-0000-0000-000000000000" {
		return nil, errors.New("session not found")
	}

	// Validate session status
	if session.Status != "waiting" {
		return nil, fmt.Errorf("cannot join session with status: %s", session.Status)
	}

	// Check if user already joined
	existingParticipant, err := dao.SessionParticipants.Ctx(ctx).
		Where("session_id", session.Id).
		Where("user_id", uid).
		One()
	if err == nil && !existingParticipant.IsEmpty() {
		return nil, errors.New("you have already joined this session")
	}

	// Check max participants limit
	participantCount, err := dao.SessionParticipants.Ctx(ctx).
		Where("session_id", session.Id).
		Count()
	if err != nil {
		return nil, errors.New("failed to check participants")
	}

	if participantCount >= session.MaxParticipants {
		return nil, errors.New("session is full")
	}

	// Create participant record
	_, err = dao.SessionParticipants.Ctx(ctx).Data(do.SessionParticipants{
		SessionId:        session.Id,
		UserId:           uid,
		Score:            0,
		Rank:             nil,
		CorrectAnswers:   0,
		WrongAnswers:     0,
		TotalTimeSeconds: 0,
		IsActive:         true,
	}).Insert()
	if err != nil {
		s.logger.Error(ctx, "Failed to create participant", "error", err)
		return nil, errors.New("failed to join session")
	}

	// Get participant ID
	var participant entity.SessionParticipants
	err = dao.SessionParticipants.Ctx(ctx).
		Where("session_id", session.Id).
		Where("user_id", uid).
		Scan(&participant)
	if err != nil {
		return nil, errors.New("failed to retrieve participant info")
	}

	// Get user info for Redis and event payload
	var user entity.Users
	err = dao.Users.Ctx(ctx).Where("id", uid).Scan(&user)
	if err != nil {
		s.logger.Error(ctx, "Failed to get user info", "error", err)
		return nil, errors.New("failed to retrieve user information")
	}

	// Update Redis participants count
	err = s.updateParticipantsCount(ctx, session.Id.String(), participantCount+1)
	if err != nil {
		s.logger.Warning(ctx, "Failed to update Redis participants count", "error", err)
	}

	// Update Redis Hash: session:{session_id}:participants
	err = s.addParticipantToRedis(ctx, session.Id.String(), uid.String(), map[string]interface{}{
		"username":  user.Username,
		"score":     0,
		"rank":      nil,
		"joined_at": participant.JoinedAt.String(),
	})
	if err != nil {
		s.logger.Warning(ctx, "Failed to add participant to Redis hash", "error", err)
	}

	// Add user to Redis Set: user:{user_id}:active_sessions
	err = s.addUserToActiveSessions(ctx, uid.String(), session.Id.String())
	if err != nil {
		s.logger.Warning(ctx, "Failed to add user to active sessions set", "error", err)
	}

	// Publish participant_joined event to Redis Pub/Sub via Centrifugo service
	err = s.centrifugoService.PublishSessionEvent(ctx, session.Id.String(), "participant_joined", map[string]interface{}{
		"user_id":            uid.String(),
		"username":           user.Username,
		"full_name":          user.FullName,
		"avatar_url":         user.AvatarUrl,
		"participants_count": participantCount + 1,
		"joined_at":          participant.JoinedAt.String(),
	})
	if err != nil {
		s.logger.Warning(ctx, "Failed to publish participant_joined event", "error", err)
	}

	// Get quiz info
	var quiz entity.Quizzes
	dao.Quizzes.Ctx(ctx).Where("id", session.QuizId).Scan(&quiz)

	// Get all current participants
	participants, err := s.getSessionParticipants(ctx, session.Id)
	if err != nil {
		s.logger.Warning(ctx, "Failed to get participants list", "error", err)
		participants = []api.ParticipantInfo{} // Return empty list on error
	}

	// Generate Centrifugo JWT token
	centrifugoToken, err := s.generateCentrifugoToken(ctx, uid.String(), session.Id.String())
	if err != nil {
		s.logger.Warning(ctx, "Failed to generate Centrifugo token", "error", err)
		centrifugoToken = "" // Continue without token
	}

	s.logger.Info(ctx, "User joined session successfully", "userId", userId, "sessionId", session.Id.String())

	return &api.JoinSessionRes{
		SessionId:     session.Id.String(),
		SessionCode:   session.SessionCode,
		QuizTitle:     quiz.Title,
		Status:        session.Status,
		ParticipantId: participant.Id,
		Participants:  participants,
		CentrifugoInfo: &api.CentrifugoAuthInfo{
			Token: centrifugoToken,
			Channels: []string{
				fmt.Sprintf("sessions#%s", session.Id.String()),
				fmt.Sprintf("leaderboard#%s", session.Id.String()),
				fmt.Sprintf("personal#%s", userId),
			},
		},
	}, nil
}

func (s *SessionServiceImpl) GetSession(ctx context.Context, sessionId string) (*api.GetSessionRes, error) {
	// Try to get from Redis cache first
	cachedInfo, err := s.getSessionFromCache(ctx, sessionId)
	if err == nil && cachedInfo != nil {
		s.logger.Info(ctx, "Session info retrieved from cache", "sessionId", sessionId)
		return cachedInfo, nil
	}

	// Cache miss, query database
	sid, err := uuid.Parse(sessionId)
	if err != nil {
		return nil, errors.New("invalid session ID")
	}

	var session entity.QuizSessions
	err = dao.QuizSessions.Ctx(ctx).Where("id", sid).Scan(&session)
	if err != nil || session.Id.String() == "00000000-0000-0000-0000-000000000000" {
		return nil, errors.New("session not found")
	}

	// Get quiz info
	var quiz entity.Quizzes
	dao.Quizzes.Ctx(ctx).Where("id", session.QuizId).Scan(&quiz)

	// Get host info
	var host entity.Users
	dao.Users.Ctx(ctx).Where("id", session.HostId).Scan(&host)

	// Get participants count
	participantsCount, _ := dao.SessionParticipants.Ctx(ctx).
		Where("session_id", sid).
		Count()

	// Calculate duration
	var durationSeconds int
	if session.StartedAt != nil && session.EndedAt != nil {
		durationSeconds = int(session.EndedAt.Time.Sub(session.StartedAt.Time).Seconds())
	} else if session.StartedAt != nil {
		durationSeconds = int(time.Since(session.StartedAt.Time).Seconds())
	}

	result := &api.GetSessionRes{
		SessionId:            session.Id.String(),
		SessionCode:          session.SessionCode,
		QuizId:               session.QuizId.String(),
		QuizTitle:            quiz.Title,
		Status:               session.Status,
		ParticipantsCount:    participantsCount,
		CurrentQuestionIndex: session.CurrentQuestionIndex,
		TotalQuestions:       quiz.TotalQuestions,
		StartedAt:            "",
		EndedAt:              "",
		DurationSeconds:      durationSeconds,
		HostId:               session.HostId.String(),
		HostUsername:         host.Username,
	}

	if session.StartedAt != nil {
		result.StartedAt = session.StartedAt.String()
	}
	if session.EndedAt != nil {
		result.EndedAt = session.EndedAt.String()
	}

	// Cache the result
	s.cacheSessionDetails(ctx, sessionId, result)

	return result, nil
}

// Helper functions

func (s *SessionServiceImpl) generateUniqueSessionCode(ctx context.Context) (string, error) {
	maxAttempts := 10
	for i := 0; i < maxAttempts; i++ {
		code := utility.GenerateSessionCode()

		// Check if code already exists in active sessions
		count, err := dao.QuizSessions.Ctx(ctx).
			Where("session_code", code).
			WhereIn("status", []string{"waiting", "in_progress"}).
			Count()

		if err != nil {
			return "", err
		}

		if count == 0 {
			return code, nil
		}
	}

	return "", errors.New("failed to generate unique session code after multiple attempts")
}

func (s *SessionServiceImpl) cacheSessionInfo(ctx context.Context, sessionId string, info map[string]interface{}) error {
	if s.redisClient == nil {
		return errors.New("Redis client not available")
	}

	key := fmt.Sprintf("session:%s:info", sessionId)

	// Set hash fields
	for field, value := range info {
		err := s.redisClient.HSet(ctx, key, field, value).Err()
		if err != nil {
			return err
		}
	}

	// Set TTL to 1 hour
	s.redisClient.Expire(ctx, key, 1*time.Hour)

	return nil
}

func (s *SessionServiceImpl) updateParticipantsCount(ctx context.Context, sessionId string, count int) error {
	if s.redisClient == nil {
		return nil
	}

	key := fmt.Sprintf("session:%s:info", sessionId)
	return s.redisClient.HSet(ctx, key, "participants_count", count).Err()
}

func (s *SessionServiceImpl) getSessionParticipants(ctx context.Context, sessionId uuid.UUID) ([]api.ParticipantInfo, error) {
	type ParticipantWithUser struct {
		entity.SessionParticipants
		Username  string
		FullName  string
		AvatarUrl string
	}

	var participants []ParticipantWithUser
	err := dao.SessionParticipants.Ctx(ctx).
		Where("sp.session_id", sessionId).
		LeftJoin("users u", "sp.user_id = u.id").
		Fields("sp.*, u.username, u.full_name, u.avatar_url").
		Scan(&participants)

	if err != nil {
		return nil, err
	}

	result := make([]api.ParticipantInfo, len(participants))
	for i, p := range participants {
		result[i] = api.ParticipantInfo{
			UserId:    p.UserId.String(),
			Username:  p.Username,
			FullName:  p.FullName,
			AvatarUrl: p.AvatarUrl,
			Score:     p.Score,
			Rank:      p.Rank,
			JoinedAt:  p.JoinedAt.String(),
		}
	}

	return result, nil
}

func (s *SessionServiceImpl) getSessionFromCache(ctx context.Context, sessionId string) (*api.GetSessionRes, error) {
	if s.redisClient == nil {
		return nil, errors.New("Redis not available")
	}

	key := fmt.Sprintf("session:%s:info", sessionId)

	// Check if key exists
	exists, err := s.redisClient.Exists(ctx, key).Result()
	if err != nil || exists == 0 {
		return nil, errors.New("cache miss")
	}

	// This is a basic cache - full implementation would reconstruct the complete response
	// For now, return nil to fall back to database query
	return nil, errors.New("cache reconstruction not implemented")
}

func (s *SessionServiceImpl) cacheSessionDetails(ctx context.Context, sessionId string, session *api.GetSessionRes) error {
	if s.redisClient == nil {
		return nil
	}

	key := fmt.Sprintf("session:%s:details", sessionId)

	// Serialize to JSON
	data, err := json.Marshal(session)
	if err != nil {
		return err
	}

	// Cache with 1 hour TTL
	return s.redisClient.Set(ctx, key, data, 1*time.Hour).Err()
}

// addParticipantToRedis adds participant info to Redis Hash
func (s *SessionServiceImpl) addParticipantToRedis(ctx context.Context, sessionId, userId string, info map[string]interface{}) error {
	if s.redisClient == nil {
		return errors.New("Redis client not available")
	}

	key := fmt.Sprintf("session:%s:participants", sessionId)

	// Serialize participant info to JSON
	data, err := json.Marshal(info)
	if err != nil {
		return err
	}

	// Store as Hash field: user_id → JSON data
	err = s.redisClient.HSet(ctx, key, userId, string(data)).Err()
	if err != nil {
		return err
	}

	// Set TTL to 1 hour
	s.redisClient.Expire(ctx, key, 1*time.Hour)

	return nil
}

// addUserToActiveSessions adds session to user's active sessions set
func (s *SessionServiceImpl) addUserToActiveSessions(ctx context.Context, userId, sessionId string) error {
	if s.redisClient == nil {
		return errors.New("Redis client not available")
	}

	key := fmt.Sprintf("user:%s:active_sessions", userId)

	// Add session ID to set
	err := s.redisClient.SAdd(ctx, key, sessionId).Err()
	if err != nil {
		return err
	}

	// Set TTL to 24 hours
	s.redisClient.Expire(ctx, key, 24*time.Hour)

	return nil
}

// generateCentrifugoToken generates a JWT token for Centrifugo authentication
func (s *SessionServiceImpl) generateCentrifugoToken(ctx context.Context, userId, sessionId string) (string, error) {
	channels := []string{
		fmt.Sprintf("sessions#%s", sessionId),
		fmt.Sprintf("leaderboard#%s", sessionId),
		fmt.Sprintf("personal#%s", userId),
	}

	token, err := s.centrifugoService.GenerateConnectionToken(ctx, userId, channels)
	if err != nil {
		s.logger.Error(ctx, "Failed to generate Centrifugo token", "error", err)
		return "", err
	}

	s.logger.Info(ctx, "Generated Centrifugo token", "userId", userId, "sessionId", sessionId)
	return token, nil
}
