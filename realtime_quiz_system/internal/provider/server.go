package provider

import (
	"context"
	"realtime_quiz_system/internal/config"
	"realtime_quiz_system/internal/controller"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/glog"
	"go.uber.org/fx"
)

// ServerParams defines dependencies for HTTP server
type ServerParams struct {
	fx.In

	Logger         *glog.Logger
	UserController *controller.UserController
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

					s.Group("/", func(group *ghttp.RouterGroup) {
						group.Middleware(CustomResponseHandler)

						group.GET("/swagger", func(r *ghttp.Request) {
							r.Response.ServeFile("resource/public/html/swagger.html")
						})

						group.Bind(params.UserController)
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
