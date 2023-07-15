package service

import (
	"context"
	"github.com/kirychukyurii/wasker/internal/app/directory/model"
	"github.com/kirychukyurii/wasker/internal/app/directory/repository"
	"github.com/kirychukyurii/wasker/internal/pkg/log"
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
	u, err := a.userRepository.ReadUser(ctx, userId)
	if err != nil {
		return nil, err
	}

	return u, nil
}
