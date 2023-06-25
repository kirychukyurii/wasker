package service

import (
	"context"
	model2 "github.com/kirychukyurii/wasker/internal/directory/model"
	repository2 "github.com/kirychukyurii/wasker/internal/directory/repository"
	"github.com/kirychukyurii/wasker/internal/model"
	"github.com/kirychukyurii/wasker/internal/pkg/server/interceptor"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/kirychukyurii/wasker/internal/errors"
	"github.com/kirychukyurii/wasker/internal/pkg/log"
)

type AuthService struct {
	logger         log.Logger
	authRepository repository2.AuthRepository
	userRepository repository2.UserRepository
}

func NewAuthService(logger log.Logger, authRepository repository2.AuthRepository, userRepository repository2.UserRepository) AuthService {
	return AuthService{
		logger:         logger,
		authRepository: authRepository,
		userRepository: userRepository,
	}
}

func (a AuthService) Login(ctx context.Context, login *model2.UserLogin) error {
	var user *model2.User

	users, err := a.userRepository.Query(ctx, &model2.UserQueryParam{UserName: login.Username})
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
	token := &model2.UserSession{
		User: model.LookupEntity{
			Id: user.Id,
		},
		NetworkIp:   interceptor.FromContext(ctx, interceptor.ClientIPCtxKey{}),
		AccessToken: model2.NewID(),
		CreatedAt:   createdAt,
		ExpiresAt:   createdAt.Add(10 * time.Hour),
	}

	err = a.authRepository.Login(ctx, token)
	if err != nil {
		return err
	}

	return nil
}

func (a AuthService) VerifyToken(ctx context.Context, token string) (*model2.UserSession, error) {
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
