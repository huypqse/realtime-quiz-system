package controller

import (
	"context"
	"realtime_quiz_system/api"
	"realtime_quiz_system/internal/service"

	"github.com/gogf/gf/v2/os/glog"
)

type QuizController struct {
	quizService service.QuizService
	logger      *glog.Logger
}

func NewQuizController(
	quizService service.QuizService,
	logger *glog.Logger,
) *QuizController {
	return &QuizController{
		quizService: quizService,
		logger:      logger,
	}
}

func (qc *QuizController) CreateQuiz(ctx context.Context, req *api.CreateQuizReq) (res *api.CreateQuizRes, err error) {
	// Get user ID from context (set by auth middleware)
	userId := ctx.Value("user_id")
	if userId == nil {
		qc.logger.Error(ctx, "User ID not found in context")
		return nil, err
	}

	qc.logger.Info(ctx, "Creating quiz", "userId", userId, "title", req.Title)

	res, err = qc.quizService.CreateQuiz(ctx, userId.(string), req)
	if err != nil {
		qc.logger.Error(ctx, "Failed to create quiz", "error", err)
		return nil, err
	}

	qc.logger.Info(ctx, "Quiz created successfully", "quizId", res.Id)
	return res, nil
}

func (qc *QuizController) ListQuizzes(ctx context.Context, req *api.ListQuizzesReq) (res *api.ListQuizzesRes, err error) {
	qc.logger.Info(ctx, "Listing quizzes", "page", req.Page, "size", req.Size)

	res, err = qc.quizService.ListQuizzes(ctx, req)
	if err != nil {
		qc.logger.Error(ctx, "Failed to list quizzes", "error", err)
		return nil, err
	}

	qc.logger.Info(ctx, "Quizzes listed successfully", "total", res.Total)
	return res, nil
}

func (qc *QuizController) GetQuiz(ctx context.Context, req *api.GetQuizReq) (res *api.GetQuizRes, err error) {
	qc.logger.Info(ctx, "Getting quiz", "quizId", req.Id)

	res, err = qc.quizService.GetQuiz(ctx, req.Id)
	if err != nil {
		qc.logger.Error(ctx, "Failed to get quiz", "error", err)
		return nil, err
	}

	qc.logger.Info(ctx, "Quiz retrieved successfully", "quizId", req.Id)
	return res, nil
}

func (qc *QuizController) UpdateQuiz(ctx context.Context, req *api.UpdateQuizReq) (res *api.UpdateQuizRes, err error) {
	// Get user ID from context (set by auth middleware)
	userId := ctx.Value("user_id")
	if userId == nil {
		qc.logger.Error(ctx, "User ID not found in context")
		return nil, err
	}

	qc.logger.Info(ctx, "Updating quiz", "userId", userId, "quizId", req.Id)

	err = qc.quizService.UpdateQuiz(ctx, userId.(string), req)
	if err != nil {
		qc.logger.Error(ctx, "Failed to update quiz", "error", err)
		return nil, err
	}

	qc.logger.Info(ctx, "Quiz updated successfully", "quizId", req.Id)
	return &api.UpdateQuizRes{Message: "Quiz updated successfully"}, nil
}

func (qc *QuizController) DeleteQuiz(ctx context.Context, req *api.DeleteQuizReq) (res *api.DeleteQuizRes, err error) {
	// Get user ID from context (set by auth middleware)
	userId := ctx.Value("user_id")
	if userId == nil {
		qc.logger.Error(ctx, "User ID not found in context")
		return nil, err
	}

	qc.logger.Info(ctx, "Deleting quiz", "userId", userId, "quizId", req.Id)

	err = qc.quizService.DeleteQuiz(ctx, userId.(string), req.Id)
	if err != nil {
		qc.logger.Error(ctx, "Failed to delete quiz", "error", err)
		return nil, err
	}

	qc.logger.Info(ctx, "Quiz deleted successfully", "quizId", req.Id)
	return &api.DeleteQuizRes{Message: "Quiz deleted successfully"}, nil
}
