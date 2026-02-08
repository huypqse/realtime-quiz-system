package logic

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"realtime_quiz_system/internal/config"
	"realtime_quiz_system/internal/service"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/redis/go-redis/v9"
)

type centrifugoLogic struct {
	cfg         *config.Config
	redisClient *redis.Client
}

func NewCentrifugoService(cfg *config.Config, redisClient *redis.Client) service.CentrifugoService {
	return &centrifugoLogic{
		cfg:         cfg,
		redisClient: redisClient,
	}
}

// GenerateConnectionToken generates a JWT token for Centrifugo connection
func (l *centrifugoLogic) GenerateConnectionToken(ctx context.Context, userID string, channels []string) (string, error) {
	ttl, err := strconv.ParseInt(l.cfg.Centrifugo.TokenTTL, 10, 64)
	if err != nil {
		ttl = 3600 // Default 1 hour
	}

	now := time.Now().Unix()
	claims := service.CentrifugoTokenClaims{
		UserID:   userID,
		Channels: channels,
		Iat:      now,
		Exp:      now + ttl,
	}

	// Create token string
	token, err := l.generateJWT(claims)
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}

	return token, nil
}

// generateJWT creates a JWT token using HMAC-SHA256
func (l *centrifugoLogic) generateJWT(claims service.CentrifugoTokenClaims) (string, error) {
	// Create header
	header := map[string]interface{}{
		"typ": "JWT",
		"alg": "HS256",
	}

	// Encode header
	headerJSON, err := json.Marshal(header)
	if err != nil {
		return "", err
	}
	headerEncoded := base64.RawURLEncoding.EncodeToString(headerJSON)

	// Encode claims
	claimsJSON, err := json.Marshal(claims)
	if err != nil {
		return "", err
	}
	claimsEncoded := base64.RawURLEncoding.EncodeToString(claimsJSON)

	// Create signature
	message := headerEncoded + "." + claimsEncoded
	signature := l.generateHMAC(message)

	// Combine parts
	token := message + "." + signature

	return token, nil
}

// generateHMAC creates HMAC-SHA256 signature
func (l *centrifugoLogic) generateHMAC(message string) string {
	h := hmac.New(sha256.New, []byte(l.cfg.Centrifugo.Secret))
	h.Write([]byte(message))
	return base64.RawURLEncoding.EncodeToString(h.Sum(nil))
}

// PublishSessionEvent publishes event to session channel (quiz:session:{session_id}:events)
func (l *centrifugoLogic) PublishSessionEvent(ctx context.Context, sessionID string, eventType string, payload interface{}) error {
	if l.redisClient == nil {
		g.Log().Warning(ctx, "Redis client is nil, skipping session event publish")
		return nil
	}

	channel := fmt.Sprintf("quiz:session:%s:events", sessionID)
	event := service.SessionEvent{
		Type:      eventType,
		Data:      payload,
		Timestamp: time.Now(),
	}

	eventJSON, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	err = l.redisClient.Publish(ctx, channel, eventJSON).Err()
	if err != nil {
		return fmt.Errorf("failed to publish to Redis: %w", err)
	}

	g.Log().Infof(ctx, "Published session event: %s to channel: %s", eventType, channel)
	return nil
}

// PublishLeaderboardUpdate publishes leaderboard update
func (l *centrifugoLogic) PublishLeaderboardUpdate(ctx context.Context, sessionID string, payload interface{}) error {
	if l.redisClient == nil {
		g.Log().Warning(ctx, "Redis client is nil, skipping leaderboard update")
		return nil
	}

	channel := fmt.Sprintf("quiz:session:%s:leaderboard", sessionID)
	event := service.SessionEvent{
		Type:      "leaderboard_update",
		Data:      payload,
		Timestamp: time.Now(),
	}

	eventJSON, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	err = l.redisClient.Publish(ctx, channel, eventJSON).Err()
	if err != nil {
		return fmt.Errorf("failed to publish to Redis: %w", err)
	}

	g.Log().Infof(ctx, "Published leaderboard update to channel: %s", channel)
	return nil
}

// PublishPersonalNotification publishes personal notification to user
func (l *centrifugoLogic) PublishPersonalNotification(ctx context.Context, userID string, payload interface{}) error {
	if l.redisClient == nil {
		g.Log().Warning(ctx, "Redis client is nil, skipping personal notification")
		return nil
	}

	channel := fmt.Sprintf("quiz:user:%s:personal", userID)
	event := service.SessionEvent{
		Type:      "personal_notification",
		Data:      payload,
		Timestamp: time.Now(),
	}

	eventJSON, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	err = l.redisClient.Publish(ctx, channel, eventJSON).Err()
	if err != nil {
		return fmt.Errorf("failed to publish to Redis: %w", err)
	}

	g.Log().Infof(ctx, "Published personal notification to channel: %s", channel)
	return nil
}
