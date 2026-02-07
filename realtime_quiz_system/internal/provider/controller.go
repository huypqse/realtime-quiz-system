package provider

import (
	"realtime_quiz_system/internal/controller"

	"go.uber.org/fx"
)

func ProvideController() fx.Option {
	return fx.Options(
		fx.Provide(controller.NewUserController),
	)
}
