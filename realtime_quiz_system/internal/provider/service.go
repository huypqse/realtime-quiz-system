package provider

import (
	"realtime_quiz_system/internal/logic"
	"realtime_quiz_system/internal/service"

	"go.uber.org/fx"
)

// ProvideServices provides all service instances
func ProvideServices() fx.Option {
	return fx.Options(
		fx.Provide(service.NewTokenService),
		fx.Provide(logic.NewUserService), // Register UserService implementation
		fx.Provide(logic.NewQuizService), // Register QuizService implementation
	)
}
