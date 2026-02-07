package controller

import (
	"context"
	user "realtime_quiz_system/api"
	"realtime_quiz_system/internal/model/do"
	"realtime_quiz_system/internal/service"

	"github.com/gogf/gf/v2/os/glog"
)

type UserController struct {
	userService service.UserService
	logger      *glog.Logger
}

func NewUserController(
	userService service.UserService,
	logger *glog.Logger,
) *UserController {
	return &UserController{
		userService: userService,
		logger:      logger,
	}
}

func (uc *UserController) Register(ctx context.Context, req *user.UserRegisterReq) (res *user.UserRegisterRes, err error) {
	uc.logger.Info(ctx, "Received user registration request",
		"username", req.Username,
		"email", req.Email)

	doReq := &do.UserRegisterReq{
		Username: req.Username,
		Email:    req.Email,
	}

	doRes, err := uc.userService.Register(ctx, doReq)
	if err != nil {
		uc.logger.Error(ctx, "User registration failed",
			"username", req.Username,
			"email", req.Email,
			"error", err)
		return nil, err
	}

	res = &user.UserRegisterRes{
		Username: doRes.Username,
		Email:    doRes.Email,
	}

	uc.logger.Info(ctx, "User registration successful",
		"username", res.Username,
		"email", res.Email)

	return res, nil
}
