package pkg

import (
	"github.com/kirychukyurii/wasker/internal/pkg/logger"
	"go.uber.org/fx"
)

// Module exports dependency
var Module = fx.Options(
	fx.Provide(logger.New),
)
