package logic

import (
	"context"
	"errors"
	"realtime_quiz_system/internal/dao"
	"realtime_quiz_system/internal/model/do"
	"realtime_quiz_system/internal/model/entity"
	"realtime_quiz_system/internal/service"
	"realtime_quiz_system/utility"

	"github.com/gogf/gf/v2/os/glog"
	"github.com/google/uuid"
)

type UserServiceImpl struct {
	logger       *glog.Logger
	tokenService service.TokenService
}

func NewUserService(
	logger *glog.Logger,
	tokenService service.TokenService,
) service.UserService {
	return &UserServiceImpl{
		logger:       logger,
		tokenService: tokenService,
	}
}

func (l *UserServiceImpl) Register(ctx context.Context, req *do.UserRegisterReq) (user *do.UserRegisterRes, err error) {
	if req.Email == "" {
		l.logger.Error(ctx, "Invalid register request: missing email",
			"email", req.Email)
		return nil, errors.New("email is required")
	}

	if req.Username == "" {
		l.logger.Error(ctx, "Invalid register request: missing username")
		return nil, errors.New("username is required")
	}

	if req.Password == "" {
		l.logger.Error(ctx, "Invalid register request: missing password")
		return nil, errors.New("password is required")
	}

	l.logger.Info(ctx, "User registration started",
		"username", req.Username,
		"email", req.Email)

	existingUser, err := dao.Users.Ctx(ctx).Where("email", req.Email).One()
	if err == nil && !existingUser.IsEmpty() {
		l.logger.Warning(ctx, "User already exists", "email", req.Email)
		return nil, errors.New("user with this email already exists")
	}

	// Hash password
	passwordHash, err := utility.HashPassword(req.Password)
	if err != nil {
		l.logger.Error(ctx, "Failed to hash password", "error", err)
		return nil, errors.New("failed to process password")
	}

	userId, err := uuid.NewV7()
	if err != nil {
		l.logger.Error(ctx, "Failed to generate UUID", "error", err)
		return nil, errors.New("failed to generate user ID")
	}

	// Insert new user
	_, err = dao.Users.Ctx(ctx).Data(do.Users{
		Id:           userId,
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: passwordHash,
	}).Insert()
	if err != nil {
		l.logger.Error(ctx, "Failed to insert user", "error", err)
		return nil, errors.New("failed to create user")
	}

	// Generate JWT token
	token, err := l.tokenService.GenerateAuthToken(userId.String(), req.Username)
	if err != nil {
		l.logger.Error(ctx, "Failed to generate token", "error", err)
		return nil, errors.New("failed to generate authentication token")
	}

	l.logger.Info(ctx, "User registered successfully",
		"username", req.Username,
		"id", userId.String())

	return &do.UserRegisterRes{
		Id:       userId.String(),
		Username: req.Username,
		Email:    req.Email,
		Token:    token,
	}, nil
}

func (l *UserServiceImpl) Login(ctx context.Context, req *do.UserLoginReq) (user *do.UserLoginRes, err error) {
	if req.UsernameOrEmail == "" {
		return nil, errors.New("username or email is required")
	}

	if req.Password == "" {
		return nil, errors.New("password is required")
	}

	l.logger.Info(ctx, "User login attempt", "usernameOrEmail", req.UsernameOrEmail)

	// Find user by username or email
	var userEntity entity.Users
	err = dao.Users.Ctx(ctx).
		Where("username = ? OR email = ?", req.UsernameOrEmail, req.UsernameOrEmail).
		Scan(&userEntity)

	if err != nil {
		l.logger.Warning(ctx, "User not found or database error", "usernameOrEmail", req.UsernameOrEmail, "error", err)
		return nil, errors.New("invalid username/email or password")
	}

	// Verify password
	err = utility.VerifyPassword(req.Password, userEntity.PasswordHash)
	if err != nil {
		l.logger.Warning(ctx, "Invalid password", "usernameOrEmail", req.UsernameOrEmail)
		return nil, errors.New("invalid username/email or password")
	}

	// Generate JWT token
	token, err := l.tokenService.GenerateAuthToken(userEntity.Id.String(), userEntity.Username)
	if err != nil {
		l.logger.Error(ctx, "Failed to generate token", "error", err)
		return nil, errors.New("failed to generate authentication token")
	}

	l.logger.Info(ctx, "User logged in successfully", "username", userEntity.Username)

	return &do.UserLoginRes{
		Id:        userEntity.Id.String(),
		Username:  userEntity.Username,
		Email:     userEntity.Email,
		FullName:  userEntity.FullName,
		AvatarUrl: userEntity.AvatarUrl,
		Token:     token,
	}, nil
}

func (l *UserServiceImpl) GetProfile(ctx context.Context, userId string) (user *do.UserProfileRes, err error) {
	// Parse UUID
	uid, err := uuid.Parse(userId)
	if err != nil {
		return nil, errors.New("invalid user ID")
	}

	// Get user info
	var userEntity entity.Users
	err = dao.Users.Ctx(ctx).Where("id", uid).Scan(&userEntity)
	if err != nil {
		l.logger.Error(ctx, "Failed to get user", "error", err)
		return nil, errors.New("user not found")
	}

	// Get user statistics
	stats, err := l.getUserStatistics(ctx, uid)
	if err != nil {
		l.logger.Warning(ctx, "Failed to get user statistics", "error", err)
		// Continue with zero stats if statistics query fails
	}

	return &do.UserProfileRes{
		Id:              userEntity.Id.String(),
		Username:        userEntity.Username,
		Email:           userEntity.Email,
		FullName:        userEntity.FullName,
		AvatarUrl:       userEntity.AvatarUrl,
		SessionsPlayed:  stats.SessionsPlayed,
		AverageScore:    stats.AverageScore,
		HighestScore:    stats.HighestScore,
		FirstPlaceCount: stats.FirstPlaceCount,
	}, nil
}

func (l *UserServiceImpl) UpdateProfile(ctx context.Context, userId string, req *do.UserUpdateProfileReq) error {
	// Parse UUID
	uid, err := uuid.Parse(userId)
	if err != nil {
		return errors.New("invalid user ID")
	}

	// Update only allowed fields
	updateData := do.Users{}
	if req.FullName != "" {
		updateData.FullName = req.FullName
	}
	if req.AvatarUrl != "" {
		updateData.AvatarUrl = req.AvatarUrl
	}

	_, err = dao.Users.Ctx(ctx).Where("id", uid).Data(updateData).Update()
	if err != nil {
		l.logger.Error(ctx, "Failed to update user profile", "error", err)
		return errors.New("failed to update profile")
	}

	l.logger.Info(ctx, "User profile updated successfully", "userId", userId)
	return nil
}

type userStatistics struct {
	SessionsPlayed  int
	AverageScore    int
	HighestScore    int
	FirstPlaceCount int
}

func (l *UserServiceImpl) getUserStatistics(ctx context.Context, userId uuid.UUID) (*userStatistics, error) {
	stats := &userStatistics{}

	// Get sessions played and scores
	type statsResult struct {
		SessionsPlayed int `json:"sessions_played"`
		AverageScore   int `json:"average_score"`
		HighestScore   int `json:"highest_score"`
	}

	var result statsResult
	err := dao.SessionParticipants.Ctx(ctx).
		Fields("COUNT(*) as sessions_played, COALESCE(AVG(score)::int, 0) as average_score, COALESCE(MAX(score), 0) as highest_score").
		Where("user_id", userId).
		Scan(&result)

	if err != nil {
		return stats, err
	}

	stats.SessionsPlayed = result.SessionsPlayed
	stats.AverageScore = result.AverageScore
	stats.HighestScore = result.HighestScore

	// Get first place count (rank = 1)
	count, err := dao.SessionParticipants.Ctx(ctx).
		Where("user_id", userId).
		Where("rank", 1).
		Count()

	if err == nil {
		stats.FirstPlaceCount = count
	}

	return stats, nil
}
