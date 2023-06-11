package register

import (
	"context"
	"errors"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"

	"github.com/kirychukyurii/wasker/gen/go/directory/v1alpha1"
	"github.com/kirychukyurii/wasker/internal/controller"
)

func GrpcDirectoryEndpoints(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error {
	if err := v1alpha1.RegisterUserServiceHandlerFromEndpoint(ctx, mux, endpoint, opts); err != nil {
		return errors.Join(err)
	}

	if err := v1alpha1.RegisterRoleServiceHandlerFromEndpoint(ctx, mux, endpoint, opts); err != nil {
		return errors.Join(err)
	}

	if err := v1alpha1.RegisterScopeServiceHandlerFromEndpoint(ctx, mux, endpoint, opts); err != nil {
		return errors.Join(err)
	}

	return nil
}

func GrpcDirectoryServiceServers(s grpc.ServiceRegistrar, controller controller.Controllers) {
	v1alpha1.RegisterUserServiceServer(s, &controller.V1alpha1.User)
}
