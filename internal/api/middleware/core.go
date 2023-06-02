package middleware

import (
	"fmt"
	"github.com/kirychukyurii/wasker/internal/pkg/handler"
	"github.com/kirychukyurii/wasker/internal/pkg/logger"
	"github.com/labstack/echo/v4"
	"runtime"
)

type CoreMiddleware struct {
	handler handler.HttpHandler
	logger  logger.Logger
}

// NewCoreMiddleware creates new database transactions middleware
func NewCoreMiddleware(handler handler.HttpHandler, logger logger.Logger) CoreMiddleware {
	return CoreMiddleware{
		handler: handler,
		logger:  logger,
	}
}

func (a CoreMiddleware) core() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			defer func() {
				if r := recover(); r != nil {
					err, ok := r.(error)
					if !ok {
						err = fmt.Errorf("%v", r)
					}

					// recovery stack
					stack := make([]byte, 4<<10)
					length := runtime.Stack(stack, true)
					msg := fmt.Sprintf("[PANIC RECOVER] %v %s\n", err, stack[:length])
					a.logger.Log.Error().Msg(msg)

					ctx.Error(err)
				}
			}()

			if err := next(ctx); err != nil {
				ctx.Error(err)
			}

			return nil
		}
	}
}

func (a CoreMiddleware) Setup() {
	a.logger.Log.Info().Msg("Setting up core middleware")
	a.handler.Engine.Use(a.core())
}
