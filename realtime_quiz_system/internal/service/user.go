package service

import (
	"context"
	"realtime_quiz_system/internal/model/do"
)

type UserService interface {
	Register(ctx context.Context, req *do.UserRegisterReq) (user *do.UserRegisterRes, err error)
}
