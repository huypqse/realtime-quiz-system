package service

import (
	"context"
	"realtime_quiz_system/api"
)

type QuizService interface {
	CreateQuiz(ctx context.Context, userId string, req *api.CreateQuizReq) (*api.CreateQuizRes, error)
	ListQuizzes(ctx context.Context, req *api.ListQuizzesReq) (*api.ListQuizzesRes, error)
	GetQuiz(ctx context.Context, quizId string) (*api.GetQuizRes, error)
	UpdateQuiz(ctx context.Context, userId string, req *api.UpdateQuizReq) error
	DeleteQuiz(ctx context.Context, userId string, quizId string) error
}
