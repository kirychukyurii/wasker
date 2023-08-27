package server

import (
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/kirychukyurii/wasker/internal/app/directory/controller"
	"github.com/kirychukyurii/wasker/internal/config"
	"github.com/kirychukyurii/wasker/internal/pkg/log"
	"github.com/kirychukyurii/wasker/internal/pkg/server/interceptor"
	"github.com/kirychukyurii/wasker/internal/pkg/server/interceptor/auth"
	"github.com/kirychukyurii/wasker/internal/pkg/server/interceptor/requestid"
)

type GrpcServer struct {
	Server *grpc.Server
}

func NewGrpcServer(cfg config.Config, logger log.Logger, controller controller.Controllers) GrpcServer {
	l, opts := interceptor.NewGrpcLoggingHandler(logger)
	r := recovery.WithRecoveryHandler(interceptor.NewGrpcPanicRecoveryHandler(logger))

	// create new gRPC server
	s := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			recovery.UnaryServerInterceptor(r),
			interceptor.ErrorUnaryServerInterceptor(),
			requestid.UnaryServerInterceptor(logger),
			logging.UnaryServerInterceptor(l, opts...),
			auth.UnaryServerInterceptor(logger, controller),
			// Add any other interceptor you want.
		),
		grpc.ChainStreamInterceptor(
			recovery.StreamServerInterceptor(r),
			logging.StreamServerInterceptor(l, opts...),
			// Add any other interceptor you want.
		))

	// Register reflection service on gRPC server.
	reflection.Register(s)

	return GrpcServer{
		Server: s,
	}
}
