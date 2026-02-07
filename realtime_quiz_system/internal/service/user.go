package service

import (
	"context"
	"realtime_quiz_system/internal/model/do"
)

type UserService interface {
	Register(ctx context.Context, req *do.UserRegisterReq) (user *do.UserRegisterRes, err error)
	Login(ctx context.Context, req *do.UserLoginReq) (user *do.UserLoginRes, err error)
	GetProfile(ctx context.Context, userId string) (user *do.UserProfileRes, err error)
	UpdateProfile(ctx context.Context, userId string, req *do.UserUpdateProfileReq) error
}
