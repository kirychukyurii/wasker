package server

import (
	"go.uber.org/fx"
)

// Module exports dependency
var Module = fx.Options(
	fx.Provide(NewGrpcServer),
	fx.Provide(NewHttpServer),
)
