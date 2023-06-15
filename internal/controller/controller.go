package controller

import "go.uber.org/fx"

var Module = fx.Options(
	fx.Provide(NewUserController),
	fx.Provide(NewAuthController),
	fx.Provide(NewGroupControllers),
)

type Controllers struct {
	User UserController
	Auth AuthController
}

func NewGroupControllers(u UserController, a AuthController) Controllers {
	return Controllers{
		User: u,
		Auth: a,
	}
}
