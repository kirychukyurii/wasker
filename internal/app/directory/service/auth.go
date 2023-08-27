package service

import (
	"context"
	"github.com/kirychukyurii/wasker/internal/lib"
	"github.com/kirychukyurii/wasker/internal/pkg/server/interceptor/requestid"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/kirychukyurii/wasker/internal/app/directory/model"
	"github.com/kirychukyurii/wasker/internal/app/directory/repository"
	"github.com/kirychukyurii/wasker/internal/errors"
	lookup "github.com/kirychukyurii/wasker/internal/model"
	"github.com/kirychukyurii/wasker/internal/pkg/log"
)

type AuthService struct {
	logger         log.Logger
	authRepository repository.AuthRepository
	userRepository repository.UserRepository
}

func NewAuthService(logger log.Logger, authRepository repository.AuthRepository, userRepository repository.UserRepository) AuthService {
	return AuthService{
		logger:         logger,
		authRepository: authRepository,
		userRepository: userRepository,
	}
}

func (a AuthService) Login(ctx context.Context, login *model.UserLogin) error {
	var user *model.User

	users, err := a.userRepository.Query(ctx, &model.UserQueryParam{UserName: login.Username})
	if err != nil || len(users.List) < 1 {
		return errors.NewBadRequestError(errors.AppError{
			Message: errors.ErrAuthInvalidCredentials.Error(),
			Details: errors.AppErrorDetail{
				Err:       err,
				ErrReason: "INVALID_CREDENTIALS",
				ErrDomain: "service.auth.login",
				RequestId: lib.FromContext(ctx, requestid.XRequestIDCtxKey{}).(string),
			},
		})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(users.List[0].Password), []byte(login.Password)); err != nil {
		// bcrypt.ErrHashTooShort
		// bcrypt.ErrMismatchedHashAndPassword
		return errors.NewBadRequestError(errors.AppError{
			Message: errors.ErrAuthInvalidCredentials.Error(),
			Details: errors.AppErrorDetail{
				Err:       err,
				ErrReason: "INVALID_CREDENTIALS",
				ErrDomain: "service.auth.login",
				RequestId: lib.FromContext(ctx, requestid.XRequestIDCtxKey{}).(string),
			},
		})
	}

	createdAt := time.Now()
	token := &model.UserSession{
		User: lookup.LookupEntity{
			Id: user.Id,
		},
		NetworkIp:   "0.0.0.0", // fixme
		AccessToken: model.NewID(),
		CreatedAt:   createdAt,
		ExpiresAt:   createdAt.Add(10 * time.Hour),
	}

	if err := a.authRepository.Login(ctx, token); err != nil {
		return err
	}

	return nil
}

func (a AuthService) Authn(ctx context.Context, token string) (*model.UserSession, error) {
	s, err := a.authRepository.Authn(ctx, token)
	if err != nil {
		return nil, err
	}

	return s, nil
}

func (a AuthService) Authz(ctx context.Context, userId int64, service, method string) (bool, error) {
	endpoint, permission, err := a.authRepository.Authz(ctx, userId, service, method)
	if err != nil {
		return false, err
	}

	ok := permission & endpoint
	if ok == 0 {
		return false, errors.NewForbiddenError(errors.AppError{
			Message: errors.ErrAuthPermissionDenied.Error(),
			Details: errors.AppErrorDetail{
				Err:       errors.ErrAuthPermissionDenied,
				ErrReason: "PERMISSION_DENIED",
				ErrDomain: "service.auth.authz",
				RequestId: lib.FromContext(ctx, requestid.XRequestIDCtxKey{}).(string),
			},
		})
	}

	return true, nil
}
