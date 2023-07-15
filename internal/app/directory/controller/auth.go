package controller

import (
	"context"

	pb "github.com/kirychukyurii/wasker/gen/go/directory/v1"
	"github.com/kirychukyurii/wasker/internal/app/directory/model"
	"github.com/kirychukyurii/wasker/internal/app/directory/service"
	"github.com/kirychukyurii/wasker/internal/errors"
	"github.com/kirychukyurii/wasker/internal/lib"
	"github.com/kirychukyurii/wasker/internal/pkg/log"
	"github.com/kirychukyurii/wasker/internal/pkg/server/interceptor/requestid"
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
	if v := request.Username; v == "" {
		return nil, errors.NewBadRequestError(errors.AppError{
			Message: errors.ErrAuthInvalidCredentials.Error(),
			Details: errors.AppErrorDetail{
				ErrId:     "controller.auth.login.invalid",
				Err:       errors.ErrAuthInvalidCredentials,
				RequestId: lib.FromContext(ctx, requestid.XRequestIDCtxKey{}).(string),
			},
		})
	}

	if err := a.authService.Login(ctx, &model.UserLogin{
		Username: request.Username,
		Password: request.Password,
	}); err != nil {
		return nil, err
	}

	return nil, nil
}

func (a AuthController) Authn(ctx context.Context, token string) (int64, error) {
	s, err := a.authService.Authn(ctx, token)
	if err != nil {
		return 0, err
	}

	return s.User.Id, nil
}

func (a AuthController) Authz(ctx context.Context, userId int64, service, method string) (bool, error) {
	ok, err := a.authService.Authz(ctx, userId, service, method)
	if err != nil {
		return false, err
	}

	return ok, nil
}
