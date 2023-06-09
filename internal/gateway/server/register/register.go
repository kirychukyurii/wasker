package register

import (
	"context"
	"errors"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	v1 "github.com/kirychukyurii/wasker/gen/go/directory/v1"
	"google.golang.org/grpc"
)

func GrpcDirectoryEndpoints(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error {
	if err := v1.RegisterUserServiceHandlerFromEndpoint(ctx, mux, endpoint, opts); err != nil {
		return errors.Join(err)
	}

	if err := v1.RegisterRoleServiceHandlerFromEndpoint(ctx, mux, endpoint, opts); err != nil {
		return errors.Join(err)
	}

	if err := v1.RegisterScopeServiceHandlerFromEndpoint(ctx, mux, endpoint, opts); err != nil {
		return errors.Join(err)
	}

	if err := v1.RegisterAuthServiceHandlerFromEndpoint(ctx, mux, endpoint, opts); err != nil {
		return errors.Join(err)
	}

	return nil
}
