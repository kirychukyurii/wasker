package service

import (
	"context"
	"github.com/kirychukyurii/wasker/internal/model"
	"github.com/pkg/errors"
	"time"

	"github.com/kirychukyurii/wasker/internal/pkg/log"
	"github.com/kirychukyurii/wasker/internal/repository"
)

type AuthService struct {
	logger         log.Logger
	authRepository repository.AuthRepository
}

func NewAuthService(logger log.Logger, authRepository repository.AuthRepository) AuthService {
	return AuthService{
		logger:         logger,
		authRepository: authRepository,
	}
}

func (a AuthService) CheckToken(ctx context.Context, token string) (*model.UserSession, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	s, err := a.authRepository.CheckToken(ctx, token)
	if err != nil {
		return nil, errors.Wrap(err, "a.authRepository.CheckToken()")
	}

	return s, nil
}
