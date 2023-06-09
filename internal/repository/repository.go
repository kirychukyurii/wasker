package repository

import (
	"github.com/kirychukyurii/wasker/internal/repository/user/v1alpha1"
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(v1alpha1.NewUserRepository),
)
