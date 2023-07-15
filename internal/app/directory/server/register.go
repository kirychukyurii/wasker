package server

import (
	"google.golang.org/grpc"

	"github.com/kirychukyurii/wasker/gen/go/directory/v1"
	"github.com/kirychukyurii/wasker/internal/app/directory/controller"
)

func GrpcDirectoryServiceServers(s grpc.ServiceRegistrar, controller controller.Controllers) {
	v1.RegisterUserServiceServer(s, &controller.User)
	v1.RegisterAuthServiceServer(s, &controller.Auth)
}
