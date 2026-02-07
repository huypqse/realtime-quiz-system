package provider

import (
	"context"
	"net/http"
	"realtime_quiz_system/internal/service"
	"strings"

	"github.com/gogf/gf/v2/net/ghttp"
)

// CustomResponseHandler is a custom middleware for handling API responses
// It wraps successful responses with code: 200 instead of code: 0
func CustomResponseHandler(r *ghttp.Request) {
	r.Middleware.Next()

	// Skip if response already written
	if r.Response.BufferLength() > 0 {
		return
	}

	var (
		msg  string
		err  = r.GetError()
		res  = r.GetHandlerResponse()
		code = http.StatusOK
	)

	if err != nil {
		code = http.StatusInternalServerError
		msg = err.Error()
	} else {
		msg = "OK"
	}

	r.Response.WriteJson(ghttp.DefaultHandlerResponse{
		Code:    code,
		Message: msg,
		Data:    res,
	})
}

// AuthMiddleware validates JWT tokens on protected endpoints
func AuthMiddleware(tokenService service.TokenService) func(r *ghttp.Request) {
	return func(r *ghttp.Request) {
		// Get Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			r.Response.WriteJsonExit(ghttp.DefaultHandlerResponse{
				Code:    http.StatusUnauthorized,
				Message: "Authorization header is required",
				Data:    nil,
			})
			return
		}

		// Check if it's a Bearer token
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			r.Response.WriteJsonExit(ghttp.DefaultHandlerResponse{
				Code:    http.StatusUnauthorized,
				Message: "Invalid authorization header format. Use: Bearer <token>",
				Data:    nil,
			})
			return
		}

		token := parts[1]

		// Validate token
		claims, err := tokenService.ValidateToken(r.Context(), token)
		if err != nil {
			r.Response.WriteJsonExit(ghttp.DefaultHandlerResponse{
				Code:    http.StatusUnauthorized,
				Message: "Invalid or expired token",
				Data:    nil,
			})
			return
		}

		// Set user info in context
		ctx := context.WithValue(r.Context(), "user_id", claims.UserId)
		ctx = context.WithValue(ctx, "username", claims.Username)
		r.SetCtx(ctx)

		// Continue to next handler
		r.Middleware.Next()
	}
}
