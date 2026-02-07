package provider

import (
	"go.uber.org/fx"
)

// Module exports all provider options as a single Fx module
var Module = fx.Options(
	ProvideConfig(),
	ProvideLogger(),
	ProvideDatabase(),
	ProvideDAOs(),
	ProvideServices(),
	ProvideServer(),
	ProvideController(),
)
