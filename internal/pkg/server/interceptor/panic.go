package interceptor

import (
	"runtime/debug"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/kirychukyurii/wasker/internal/pkg/log"
)

func NewGrpcPanicRecoveryHandler(logger log.Logger) func(any) error {
	return grpcPanicRecoveryHandler(logger)
}

func grpcPanicRecoveryHandler(logger log.Logger) func(any) error {
	return func(p any) (err error) {
		logger.Log.Error().Err(err).Msgf("recovered from panic: %s", debug.Stack())
		return status.Errorf(codes.Internal, "%s", p)
	}
}
