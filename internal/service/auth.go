package service

import (
	"context"
	"github.com/kirychukyurii/wasker/internal/server/interceptor"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/kirychukyurii/wasker/internal/errors"
	"github.com/kirychukyurii/wasker/internal/model"
	"github.com/kirychukyurii/wasker/internal/model/dto"
	"github.com/kirychukyurii/wasker/internal/pkg/log"
	"github.com/kirychukyurii/wasker/internal/repository"
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
	if err != nil {
		return err
	}

	if len(users.List) > 0 {
		user = users.List[0]
	} else {
		return errors.ErrAuthIncorrectCredentials
	}

	hashedPassword := user.Password
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(login.Password)); err != nil {
		// bcrypt.ErrHashTooShort
		// bcrypt.ErrMismatchedHashAndPassword
		return errors.ErrAuthIncorrectCredentials
	}

	createdAt := time.Now()
	token := &model.UserSession{
		User: dto.LookupEntity{
			Id: user.Id,
		},
		NetworkIp:   ctx.Value(interceptor.ClientIPCtxKey).(string),
		AccessToken: model.NewID(),
		CreatedAt:   createdAt,
		ExpiresAt:   createdAt.Add(10 * time.Hour),
	}

	err = a.authRepository.Login(ctx, token)
	if err != nil {
		return err
	}

	return nil
}

func (a AuthService) VerifyToken(ctx context.Context, token string) (*model.UserSession, error) {
	s, err := a.authRepository.VerifyToken(ctx, token)
	if err != nil {
		return nil, err
	}

	return s, nil
}

func (a AuthService) VerifyPermission(ctx context.Context, userId uint64, service, method string) (bool, error) {
	endpoint, permission, err := a.authRepository.VerifyPermission(ctx, userId, service, method)
	if err != nil {
		return false, errors.ErrAuthPermissionDenied
	}

	ok := permission & endpoint
	if ok == 0 {
		return false, errors.ErrAuthPermissionDenied
	}

	return true, nil
}
