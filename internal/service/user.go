package service

import (
	"context"
	"time"

	"github.com/kirychukyurii/wasker/internal/model"
	"github.com/kirychukyurii/wasker/internal/pkg/log"
	"github.com/kirychukyurii/wasker/internal/repository"
)

type UserService struct {
	logger         log.Logger
	userRepository repository.UserRepository
}

func NewUserService(logger log.Logger, userRepository repository.UserRepository) UserService {
	return UserService{
		logger:         logger,
		userRepository: userRepository,
	}
}

func (a UserService) ReadUser(ctx context.Context, userId int64) (*model.User, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	u, err := a.userRepository.ReadUser(ctx, userId)
	if err != nil {
		return nil, err
	}

	return u, nil
}
