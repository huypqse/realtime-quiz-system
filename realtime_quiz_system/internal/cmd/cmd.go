package cmd

import (
	"realtime_quiz_system/internal/provider"

	"go.uber.org/fx"
)

// RunApp initializes and runs the application with Fx
func RunApp() {
	app := fx.New(
		// Configure Fx logger to be less verbose
		fx.NopLogger,

		// Load all providers
		provider.Module,
	)

	app.Run()
}
