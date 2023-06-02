package route

import (
	"github.com/kirychukyurii/wasker/internal/pkg/handler"
	"github.com/kirychukyurii/wasker/internal/pkg/logger"
)

type AuthRoutes struct {
	logger  logger.Logger
	handler handler.HttpHandler
}

func NewAuthRoutes(logger logger.Logger, handler handler.HttpHandler) AuthRoutes {
	return AuthRoutes{
		handler: handler,
		logger:  logger,
	}
}

func (a AuthRoutes) Setup() {
	a.logger.Log.Info().Msg("Setting up auth routes")

	//api := a.handler.RouterV1.Group("/auth")
	/*
		{
			api.GET("/info", a.authController.UserInfo)
			api.POST("/login", a.authController.UserLogin)
			api.POST("/logout", a.authController.UserLogout)
		}

	*/
}
