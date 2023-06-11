package v1alpha1

import (
	"context"
	"time"

	"github.com/kirychukyurii/wasker/internal/model"
	"github.com/kirychukyurii/wasker/internal/pkg/log"
	"github.com/kirychukyurii/wasker/internal/repository/user/v1alpha1"
)

type UserService struct {
	logger         log.Logger
	userRepository v1alpha1.UserRepository
}

func NewUserService(logger log.Logger, userRepository v1alpha1.UserRepository) UserService {
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
