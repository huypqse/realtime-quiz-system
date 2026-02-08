package controller

import (
	"context"
	"realtime_quiz_system/api"
	"realtime_quiz_system/internal/service"

	"github.com/gogf/gf/v2/os/glog"
)

type SessionController struct {
	sessionService service.SessionService
	logger         *glog.Logger
}

func NewSessionController(
	sessionService service.SessionService,
	logger *glog.Logger,
) *SessionController {
	return &SessionController{
		sessionService: sessionService,
		logger:         logger,
	}
}

func (sc *SessionController) CreateSession(ctx context.Context, req *api.CreateSessionReq) (res *api.CreateSessionRes, err error) {
	// Get user ID from context (set by auth middleware)
	userId := ctx.Value("user_id")
	if userId == nil {
		sc.logger.Error(ctx, "User ID not found in context")
		return nil, err
	}

	sc.logger.Info(ctx, "Creating session", "userId", userId, "quizId", req.QuizId)

	res, err = sc.sessionService.CreateSession(ctx, userId.(string), req)
	if err != nil {
		sc.logger.Error(ctx, "Failed to create session", "error", err)
		return nil, err
	}

	sc.logger.Info(ctx, "Session created successfully", "sessionId", res.SessionId)
	return res, nil
}

func (sc *SessionController) JoinSession(ctx context.Context, req *api.JoinSessionReq) (res *api.JoinSessionRes, err error) {
	// Get user ID from context (set by auth middleware)
	userId := ctx.Value("user_id")
	if userId == nil {
		sc.logger.Error(ctx, "User ID not found in context")
		return nil, err
	}

	sc.logger.Info(ctx, "User joining session", "userId", userId, "sessionCode", req.SessionCode)

	res, err = sc.sessionService.JoinSession(ctx, userId.(string), req)
	if err != nil {
		sc.logger.Error(ctx, "Failed to join session", "error", err)
		return nil, err
	}

	sc.logger.Info(ctx, "User joined session successfully", "sessionId", res.SessionId)
	return res, nil
}

func (sc *SessionController) GetSession(ctx context.Context, req *api.GetSessionReq) (res *api.GetSessionRes, err error) {
	sc.logger.Info(ctx, "Getting session details", "sessionId", req.SessionId)

	res, err = sc.sessionService.GetSession(ctx, req.SessionId)
	if err != nil {
		sc.logger.Error(ctx, "Failed to get session", "error", err)
		return nil, err
	}

	sc.logger.Info(ctx, "Session details retrieved successfully", "sessionId", req.SessionId)
	return res, nil
}
