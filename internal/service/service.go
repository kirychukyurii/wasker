package service

import (
	"go.uber.org/fx"

	"github.com/kirychukyurii/wasker/internal/service/user/v1alpha1"
)

var Module = fx.Options(
	fx.Provide(v1alpha1.NewUserService),
)
