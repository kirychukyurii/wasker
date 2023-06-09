package controller

import (
	"github.com/kirychukyurii/wasker/internal/controller/user/v1alpha1"
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(v1alpha1.NewUserController),
)
