package service

import (
	"context"
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
		User: lookup.LookupEntity{
			Id: user.Id,
		},
		NetworkIp:   "0.0.0.0", // fixme
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
		return false, errors.ErrAuthPermissionDenied
	}

	ok := permission & endpoint
	if ok == 0 {
		return false, errors.ErrAuthPermissionDenied
	}

	return true, nil
}
