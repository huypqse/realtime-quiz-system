package service

import (
	"context"
	"time"
)

type CentrifugoService interface {
	GenerateConnectionToken(ctx context.Context, userID string, channels []string) (string, error)
	PublishSessionEvent(ctx context.Context, sessionID string, eventType string, payload interface{}) error
	PublishLeaderboardUpdate(ctx context.Context, sessionID string, payload interface{}) error
	PublishPersonalNotification(ctx context.Context, userID string, payload interface{}) error
}

type CentrifugoTokenClaims struct {
	UserID   string   `json:"sub"`
	Channels []string `json:"channels"`
	Exp      int64    `json:"exp"`
	Iat      int64    `json:"iat"`
}

type SessionEvent struct {
	Type      string      `json:"type"`
	Data      interface{} `json:"data"`
	Timestamp time.Time   `json:"timestamp"`
}
