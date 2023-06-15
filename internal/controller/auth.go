package controller

import (
	"context"
	"github.com/kirychukyurii/wasker/internal/pkg/log"
	"github.com/kirychukyurii/wasker/internal/service"
	"github.com/pkg/errors"
)

// AuthController will implement the service defined in protocol buffer definitions
type AuthController struct {
	authService service.AuthService
	logger      log.Logger
}

func NewAuthController(authService service.AuthService, logger log.Logger) AuthController {
	return AuthController{
		authService: authService,
		logger:      logger,
	}
}

func (a AuthController) CheckToken(ctx context.Context, token string, service string, method string) error {
	_, err := a.authService.CheckToken(ctx, token)
	if err != nil {
		return errors.Wrap(err, "a.authService.CheckToken()")
	}

	return nil
}
