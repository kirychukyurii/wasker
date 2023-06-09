package service

import (
	"github.com/kirychukyurii/wasker/internal/service/user/v1alpha1"
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(v1alpha1.NewUserService),
)
