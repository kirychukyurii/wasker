package pkg

import (
	"github.com/kirychukyurii/wasker/internal/pkg/db"
	"github.com/kirychukyurii/wasker/internal/pkg/handler"
	"github.com/kirychukyurii/wasker/internal/pkg/log"
	"go.uber.org/fx"
)

// Module exports dependency
var Module = fx.Options(
	fx.Provide(log.New),
	fx.Provide(db.New),
	fx.Provide(handler.New),
)
