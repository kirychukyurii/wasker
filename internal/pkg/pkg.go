package pkg

import (
	"go.uber.org/fx"

	"github.com/kirychukyurii/wasker/internal/pkg/consul"
	"github.com/kirychukyurii/wasker/internal/pkg/db"
	"github.com/kirychukyurii/wasker/internal/pkg/log"
	"github.com/kirychukyurii/wasker/internal/pkg/server"
)

// Module exports dependency
var Module = fx.Options(
	fx.Provide(log.New),
	fx.Provide(db.New),
	fx.Provide(consul.New),
	fx.Provide(server.NewGrpcServer),
	fx.Provide(server.NewHttpServer),
)
