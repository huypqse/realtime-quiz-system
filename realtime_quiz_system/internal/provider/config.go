package provider

import (
	"context"
	"realtime_quiz_system/internal/config"

	"go.uber.org/fx"
)

// ProvideConfig provides the application configuration
func ProvideConfig() fx.Option {
	return fx.Options(
		fx.Provide(func() *config.Config {
			// Initialize config immediately when provider is called
			ctx := context.Background()
			config.InitConfig(ctx)

			// Return pointer to the config
			cfg := config.GetConfig()
			return &cfg
		}),
	)
}
