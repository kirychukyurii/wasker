package controller

import (
	"go.uber.org/fx"

	"github.com/kirychukyurii/wasker/internal/controller/v1alpha1"
)

var Module = fx.Options(
	fx.Provide(v1alpha1.NewUserController),
	fx.Provide(v1alpha1.NewV1alpha1Controllers),
	fx.Provide(NewGroupControllers),
)

type Controllers struct {
	V1alpha1 v1alpha1.V1alpha1
}

func NewGroupControllers(v1alpha1 v1alpha1.V1alpha1) Controllers {
	return Controllers{
		V1alpha1: v1alpha1,
	}
}
