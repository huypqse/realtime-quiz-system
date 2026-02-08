package provider

import (
	"context"
	"fmt"
	"realtime_quiz_system/internal/config"
	"time"

	"github.com/gogf/gf/v2/os/glog"
	"github.com/redis/go-redis/v9"
	"go.uber.org/fx"
)

// ProvideRedis provides Redis client
func ProvideRedis() fx.Option {
	return fx.Options(
		fx.Provide(func(lc fx.Lifecycle, cfg *config.Config, logger *glog.Logger) *redis.Client {
			redisClient := redis.NewClient(&redis.Options{
				Addr:            fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port),
				Password:        cfg.Redis.Password,
				DB:              cfg.Redis.DB,
				MaxIdleConns:    cfg.Redis.MaxIdle,
				MaxActiveConns:  cfg.Redis.MaxActive,
				ConnMaxIdleTime: parseDuration(cfg.Redis.IdleTimeout),
			})

			lc.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					logger.Info(ctx, "Connecting to Redis...")

					// Test connection
					_, err := redisClient.Ping(ctx).Result()
					if err != nil {
						logger.Warning(ctx, "Failed to connect to Redis", "error", err)
						logger.Warning(ctx, "Redis features will be disabled")
						return nil // Don't fail startup if Redis is unavailable
					}

					logger.Info(ctx, "Redis connected successfully")
					return nil
				},
				OnStop: func(ctx context.Context) error {
					logger.Info(ctx, "Closing Redis connection...")
					return redisClient.Close()
				},
			})

			return redisClient
		}),
	)
}

func parseDuration(s string) time.Duration {
	d, err := time.ParseDuration(s)
	if err != nil {
		return 5 * time.Minute // Default
	}
	return d
}
