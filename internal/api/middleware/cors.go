package middleware

import (
	"github.com/kirychukyurii/wasker/internal/pkg/handler"
	"github.com/kirychukyurii/wasker/internal/pkg/logger"
	"github.com/labstack/echo/v4/middleware"
)

// CorsMiddleware middleware for cors
type CorsMiddleware struct {
	handler handler.HttpHandler
	logger  logger.Logger
}

// NewCorsMiddleware creates new cors middleware
func NewCorsMiddleware(handler handler.HttpHandler, logger logger.Logger) CorsMiddleware {
	return CorsMiddleware{
		handler: handler,
		logger:  logger,
	}
}

func (a CorsMiddleware) Setup() {
	a.logger.Log.Info().Msg("Setting up cors middleware")

	a.handler.Engine.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"*"},
		AllowCredentials: true,
		AllowHeaders:     []string{"*"},
		AllowMethods:     []string{"*"},
	}))
}
