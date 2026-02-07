package logic

import (
	"context"
	"errors"
	"realtime_quiz_system/internal/dao"
	"realtime_quiz_system/internal/model/do"
	"realtime_quiz_system/internal/service"

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

	l.logger.Info(ctx, "User registration started",
		"username", req.Username,
		"email", req.Email)

	existingUser, err := dao.Users.Ctx(ctx).Where("email", req.Email).One()
	if err == nil && !existingUser.IsEmpty() {
		l.logger.Warning(ctx, "User already exists", "email", req.Email)
		return nil, errors.New("user with this email already exists")
	}

	userId, err := uuid.NewV7()
	if err != nil {
		l.logger.Error(ctx, "Failed to generate UUID", "error", err)
		return nil, errors.New("failed to generate user ID")
	}

	// Insert new user
	_, err = dao.Users.Ctx(ctx).Data(do.Users{
		Id:       userId,
		Username: req.Username,
		Email:    req.Email,
	}).Insert()
	if err != nil {
		l.logger.Error(ctx, "Failed to insert user", "error", err)
		return nil, errors.New("failed to create user")
	}

	l.logger.Info(ctx, "User registered successfully",
		"username", req.Username,
		"id", userId.String())

	return &do.UserRegisterRes{
		Username: req.Username,
		Email:    req.Email,
	}, nil
}
