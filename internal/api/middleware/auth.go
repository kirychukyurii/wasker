package middleware

import (
	"github.com/kirychukyurii/wasker/internal/config"
	"github.com/kirychukyurii/wasker/internal/pkg/handler"
	"github.com/kirychukyurii/wasker/internal/pkg/logger"
	"github.com/labstack/echo/v4"
)

type AuthMiddleware struct {
	config  config.Config
	handler handler.HttpHandler
	logger  logger.Logger
}

func NewAuthMiddleware(config config.Config, handler handler.HttpHandler, logger logger.Logger) AuthMiddleware {
	return AuthMiddleware{
		config:  config,
		handler: handler,
		logger:  logger,
	}
}

func (a AuthMiddleware) core() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			return next(ctx)
		}
	}
}

func (a AuthMiddleware) Setup() {
	a.logger.Log.Info().Msg("setting up auth middleware")
	a.handler.Engine.Use(a.core())
}
