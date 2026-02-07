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
		Password: req.Password,
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
		Id:       doRes.Id,
		Username: doRes.Username,
		Email:    doRes.Email,
		Token:    doRes.Token,
	}

	uc.logger.Info(ctx, "User registration successful",
		"username", res.Username,
		"email", res.Email)

	return res, nil
}

func (uc *UserController) Login(ctx context.Context, req *user.UserLoginReq) (res *user.UserLoginRes, err error) {
	uc.logger.Info(ctx, "Received user login request",
		"usernameOrEmail", req.UsernameOrEmail)

	doReq := &do.UserLoginReq{
		UsernameOrEmail: req.UsernameOrEmail,
		Password:        req.Password,
	}

	doRes, err := uc.userService.Login(ctx, doReq)
	if err != nil {
		uc.logger.Error(ctx, "User login failed",
			"usernameOrEmail", req.UsernameOrEmail,
			"error", err)
		return nil, err
	}

	res = &user.UserLoginRes{
		Id:        doRes.Id,
		Username:  doRes.Username,
		Email:     doRes.Email,
		FullName:  doRes.FullName,
		AvatarUrl: doRes.AvatarUrl,
		Token:     doRes.Token,
	}

	uc.logger.Info(ctx, "User login successful",
		"username", res.Username)

	return res, nil
}

func (uc *UserController) Profile(ctx context.Context, req *user.UserProfileReq) (res *user.UserProfileRes, err error) {
	// Get user ID from context (set by auth middleware)
	userId := ctx.Value("user_id")
	if userId == nil {
		uc.logger.Error(ctx, "User ID not found in context")
		return nil, err
	}

	userIdStr, ok := userId.(string)
	if !ok {
		uc.logger.Error(ctx, "Invalid user ID in context")
		return nil, err
	}

	uc.logger.Info(ctx, "Received get profile request", "userId", userIdStr)

	doRes, err := uc.userService.GetProfile(ctx, userIdStr)
	if err != nil {
		uc.logger.Error(ctx, "Failed to get user profile",
			"userId", userIdStr,
			"error", err)
		return nil, err
	}

	res = &user.UserProfileRes{
		Id:              doRes.Id,
		Username:        doRes.Username,
		Email:           doRes.Email,
		FullName:        doRes.FullName,
		AvatarUrl:       doRes.AvatarUrl,
		SessionsPlayed:  doRes.SessionsPlayed,
		AverageScore:    doRes.AverageScore,
		HighestScore:    doRes.HighestScore,
		FirstPlaceCount: doRes.FirstPlaceCount,
	}

	uc.logger.Info(ctx, "Get profile successful", "userId", userIdStr)

	return res, nil
}

func (uc *UserController) UpdateProfile(ctx context.Context, req *user.UserUpdateProfileReq) (res *user.UserUpdateProfileRes, err error) {
	// Get user ID from context (set by auth middleware)
	userId := ctx.Value("user_id")
	if userId == nil {
		uc.logger.Error(ctx, "User ID not found in context")
		return nil, err
	}

	userIdStr, ok := userId.(string)
	if !ok {
		uc.logger.Error(ctx, "Invalid user ID in context")
		return nil, err
	}

	uc.logger.Info(ctx, "Received update profile request", "userId", userIdStr)

	doReq := &do.UserUpdateProfileReq{
		FullName:  req.FullName,
		AvatarUrl: req.AvatarUrl,
	}

	err = uc.userService.UpdateProfile(ctx, userIdStr, doReq)
	if err != nil {
		uc.logger.Error(ctx, "Failed to update user profile",
			"userId", userIdStr,
			"error", err)
		return nil, err
	}

	res = &user.UserUpdateProfileRes{
		Message: "Profile updated successfully",
	}

	uc.logger.Info(ctx, "Update profile successful", "userId", userIdStr)

	return res, nil
}
