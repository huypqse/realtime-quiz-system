package provider

import (
	"context"
	"realtime_quiz_system/internal/config"
	"realtime_quiz_system/internal/controller"
	"realtime_quiz_system/internal/service"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/net/goai"
	"github.com/gogf/gf/v2/os/glog"
	"go.uber.org/fx"
)

// ServerParams defines dependencies for HTTP server
type ServerParams struct {
	fx.In

	Logger         *glog.Logger
	UserController *controller.UserController
	QuizController *controller.QuizController
	TokenService   service.TokenService
}

// ProvideServer provides the HTTP server
func ProvideServer() fx.Option {
	return fx.Options(
		fx.Invoke(func(lc fx.Lifecycle, params ServerParams) {
			s := g.Server()
			cfg := config.GetConfig()

			lc.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					params.Logger.Info(ctx, "Starting HTTP server...")

					s.SetOpenApiPath(cfg.Server.OpenAPIPath)
					s.SetSwaggerPath("")
					s.Use(ghttp.MiddlewareCORS)

					openapi := s.GetOpenApi()
					openapi.Info.Title = "Realtime Quiz System API"
					openapi.Info.Version = "1.0.0"
					openapi.Info.Description = "API documentation for the Realtime Quiz System"

					// Add Bearer token authentication to OpenAPI spec
					openapi.Components.SecuritySchemes = goai.SecuritySchemes{
						"BearerAuth": goai.SecuritySchemeRef{
							Value: &goai.SecurityScheme{
								Type:         "http",
								Scheme:       "bearer",
								BearerFormat: "JWT",
								Description:  "Enter your JWT token in the format: your-token-here",
							},
						},
					}

					s.Group("/", func(group *ghttp.RouterGroup) {
						group.Middleware(CustomResponseHandler)

						group.GET("/swagger", func(r *ghttp.Request) {
							r.Response.ServeFile("resource/public/html/swagger.html")
						})

						// Public routes (no authentication required)
						group.POST("/register", params.UserController.Register)
						group.POST("/login", params.UserController.Login)

						// Public quiz routes
						group.GET("/quizzes", params.QuizController.ListQuizzes)
						group.GET("/quizzes/{id}", params.QuizController.GetQuiz)

						// Protected routes (authentication required)
						group.Group("/", func(authGroup *ghttp.RouterGroup) {
							authGroup.Middleware(AuthMiddleware(params.TokenService))

							// User profile routes
							authGroup.GET("/profile", params.UserController.Profile)
							authGroup.PUT("/profile", params.UserController.UpdateProfile)

							// Quiz management routes (requires authentication)
							authGroup.POST("/quizzes", params.QuizController.CreateQuiz)
							authGroup.PUT("/quizzes/{id}", params.QuizController.UpdateQuiz)
							authGroup.DELETE("/quizzes/{id}", params.QuizController.DeleteQuiz)
						})
					})

					go func() {
						s.Run()
					}()

					params.Logger.Info(ctx, "HTTP server started successfully")
					return nil
				},
				OnStop: func(ctx context.Context) error {
					params.Logger.Info(ctx, "Shutting down HTTP server...")
					s.Shutdown()
					return nil
				},
			})
		}),
	)
}
