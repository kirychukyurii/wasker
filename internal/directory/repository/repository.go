package repository

import "go.uber.org/fx"

var Module = fx.Options(
	fx.Provide(NewUserRepository),
	fx.Provide(NewAuthRepository),
	fx.Provide(NewScopeRepository),
)
