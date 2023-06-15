package server

import (
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"google.golang.org/grpc"

	"github.com/kirychukyurii/wasker/internal/controller"
	"github.com/kirychukyurii/wasker/internal/pkg/log"
	"github.com/kirychukyurii/wasker/internal/server/interceptor"
	"github.com/kirychukyurii/wasker/internal/server/register"
)

type GrpcServer struct {
	Server *grpc.Server
}

func NewGrpcServer(logger log.Logger, controller controller.Controllers) GrpcServer {
	l, opts := interceptor.NewGrpcLoggingHandler(logger)
	r := recovery.WithRecoveryHandler(interceptor.NewGrpcPanicRecoveryHandler(logger))

	// create new gRPC server
	s := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			interceptor.AuthUnaryServerInterceptor(logger, controller),
			logging.UnaryServerInterceptor(l, opts...),
			recovery.UnaryServerInterceptor(r),
			// Add any other interceptor you want.
		),
		grpc.ChainStreamInterceptor(
			logging.StreamServerInterceptor(l, opts...),
			recovery.StreamServerInterceptor(r),
			// Add any other interceptor you want.
		))

	// register the UserController on the gRPC server
	register.GrpcDirectoryServiceServers(s, controller)

	return GrpcServer{
		Server: s,
	}
}
