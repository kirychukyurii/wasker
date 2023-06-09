package server

import (
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	pbUser "github.com/kirychukyurii/wasker/gen/go/user/v1alpha1"
	cUser "github.com/kirychukyurii/wasker/internal/controller/user/v1alpha1"
	"github.com/kirychukyurii/wasker/internal/pkg/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"runtime/debug"
)

type GrpcServer struct {
	Server *grpc.Server
}

func NewGrpcServer(logger log.Logger, userServiceServer cUser.UserController) GrpcServer {
	opts := []logging.Option{
		logging.WithLogOnEvents(logging.FinishCall),
		// Add any other option (check functions starting with logging.With).
	}

	grpcPanicRecoveryHandler := func(p any) (err error) {
		logger.Log.Error().Err(err).Msgf("recovered from panic: %s", debug.Stack())
		return status.Errorf(codes.Internal, "%s", p)
	}

	// create new gRPC server
	s := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			logging.UnaryServerInterceptor(log.InterceptorLogger(logger.Log), opts...),
			recovery.UnaryServerInterceptor(recovery.WithRecoveryHandler(grpcPanicRecoveryHandler)),
			// Add any other interceptor you want.
		),
		grpc.ChainStreamInterceptor(
			logging.StreamServerInterceptor(log.InterceptorLogger(logger.Log), opts...),
			recovery.StreamServerInterceptor(recovery.WithRecoveryHandler(grpcPanicRecoveryHandler)),
			// Add any other interceptor you want.
		))

	// register the GreeterServerImpl on the gRPC server
	pbUser.RegisterUserServiceServer(s, &userServiceServer)

	return GrpcServer{
		Server: s,
	}
}
