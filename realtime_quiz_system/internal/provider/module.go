package provider

import (
	"go.uber.org/fx"
)

// Module exports all provider options as a single Fx module
var Module = fx.Options(
	ProvideConfig(),
	ProvideLogger(),
	ProvideDatabase(),
	ProvideRedis(),
	ProvideDAOs(),
	ProvideServices(),
	ProvideServer(),
	ProvideController(),
)
