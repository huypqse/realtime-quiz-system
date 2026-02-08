package service

import (
	"context"
	"realtime_quiz_system/api"
)

type SessionService interface {
	CreateSession(ctx context.Context, userId string, req *api.CreateSessionReq) (*api.CreateSessionRes, error)
	JoinSession(ctx context.Context, userId string, req *api.JoinSessionReq) (*api.JoinSessionRes, error)
	GetSession(ctx context.Context, sessionId string) (*api.GetSessionRes, error)
}
