package provider

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"go.uber.org/fx"
)

// ProvideDatabase provides database connection
func ProvideDatabase() fx.Option {
	return fx.Options(
		fx.Provide(func(lc fx.Lifecycle) gdb.DB {
			db := g.DB()

			lc.Append(fx.Hook{
				OnStop: func(ctx context.Context) error {
					// Database cleanup if needed
					return nil
				},
			})

			return db
		}),
	)
}
