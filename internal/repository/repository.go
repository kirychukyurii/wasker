package repository

import (
	"go.uber.org/fx"

	"github.com/kirychukyurii/wasker/internal/repository/user/v1alpha1"
)

var Module = fx.Options(
	fx.Provide(v1alpha1.NewUserRepository),
)
