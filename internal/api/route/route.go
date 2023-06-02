package route

import "go.uber.org/fx"

// Module exports dependency to container
var Module = fx.Options(
	fx.Provide(NewAuthRoutes),
	fx.Provide(New),
)

// Routes contains multiple routes
type Routes []Route

// Route interface
type Route interface {
	Setup()
}

// New sets up routes
func New(authRoutes AuthRoutes) Routes {
	return Routes{
		authRoutes,
	}
}

func (a Routes) Setup() {
	for _, route := range a {
		route.Setup()
	}
}
