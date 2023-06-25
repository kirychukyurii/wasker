package register

import (
	"github.com/kirychukyurii/wasker/internal/directory/controller"

	"google.golang.org/grpc"

	"github.com/kirychukyurii/wasker/gen/go/directory/v1"
)

func GrpcDirectoryServiceServers(s grpc.ServiceRegistrar, controller controller.Controllers) {
	v1.RegisterUserServiceServer(s, &controller.User)
	v1.RegisterAuthServiceServer(s, &controller.Auth)
}
