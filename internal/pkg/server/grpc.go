package server

import (
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"github.com/kirychukyurii/wasker/internal/config"
	"github.com/kirychukyurii/wasker/internal/directory/controller"
	interceptor2 "github.com/kirychukyurii/wasker/internal/pkg/server/interceptor"
	"github.com/kirychukyurii/wasker/internal/pkg/server/interceptor/auth"
	"github.com/kirychukyurii/wasker/internal/pkg/server/interceptor/requestid"
	"google.golang.org/grpc"

	"github.com/kirychukyurii/wasker/internal/pkg/log"
)

type GrpcServer struct {
	Server *grpc.Server
}

func NewGrpcServer(cfg config.Config, logger log.Logger, controller controller.Controllers) GrpcServer {
	l, opts := interceptor2.NewGrpcLoggingHandler(logger)
	r := recovery.WithRecoveryHandler(interceptor2.NewGrpcPanicRecoveryHandler(logger))

	// create new gRPC server
	s := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			requestid.UnaryServerInterceptor(),
			interceptor2.ContextUnaryServerInterceptor(cfg, logger),
			auth.UnaryServerInterceptor(logger, controller),
			logging.UnaryServerInterceptor(l, opts...),
			recovery.UnaryServerInterceptor(r),
			// Add any other interceptor you want.
		),
		grpc.ChainStreamInterceptor(
			logging.StreamServerInterceptor(l, opts...),
			recovery.StreamServerInterceptor(r),
			// Add any other interceptor you want.
		))

	return GrpcServer{
		Server: s,
	}
}
