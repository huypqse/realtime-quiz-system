package provider

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/glog"
	"go.uber.org/fx"
)

// ProvideLogger provides logger instance
func ProvideLogger() fx.Option {
	return fx.Options(
		fx.Provide(func(lc fx.Lifecycle) *glog.Logger {
			logger := g.Log()

			lc.Append(fx.Hook{
				OnStart: func(context.Context) error {
					// Configure JSON handler for logging
					logger.SetHandlers(glog.HandlerJson)
					glog.SetHandlers(glog.HandlerJson)
					return nil
				},
			})

			return logger
		}),
	)
}
