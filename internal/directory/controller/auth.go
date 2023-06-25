package controller

import (
	"context"
	pb "github.com/kirychukyurii/wasker/gen/go/directory/v1"
	"github.com/kirychukyurii/wasker/internal/directory/model"
	"github.com/kirychukyurii/wasker/internal/directory/service"
	"github.com/kirychukyurii/wasker/internal/errors"
	"github.com/kirychukyurii/wasker/internal/pkg/log"
)

// AuthController will implement the service defined in protocol buffer definitions
type AuthController struct {
	pb.UnimplementedAuthServiceServer

	authService service.AuthService
	logger      log.Logger
}

func NewAuthController(authService service.AuthService, logger log.Logger) AuthController {
	return AuthController{
		authService: authService,
		logger:      logger,
	}
}

func (a AuthController) Login(ctx context.Context, request *pb.LoginRequest) (*pb.LoginResponse, error) {
	var login model.UserLogin
	if v := request.Username; v == "" {
		return nil, errors.New(errors.ErrAuthInvalidCredentials)
	}

	login.Username = request.Username
	login.Password = request.Password
	err := a.authService.Login(ctx, &login)
	if err != nil {
		return nil, errors.New(err)
	}

	return nil, nil
}

func (a AuthController) VerifyToken(ctx context.Context, token string) (uint64, error) {
	s, err := a.authService.VerifyToken(ctx, token)
	if err != nil {
		return 0, err
	}

	return s.User.Id, nil
}

func (a AuthController) VerifyPermission(ctx context.Context, userId uint64, service, method string) (bool, error) {
	ok, err := a.authService.VerifyPermission(ctx, userId, service, method)
	if err != nil {
		return false, err
	}

	return ok, nil
}
