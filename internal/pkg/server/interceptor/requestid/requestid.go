package requestid

import (
	"context"
	"github.com/kirychukyurii/wasker/internal/lib"
	"github.com/kirychukyurii/wasker/internal/pkg/log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type XRequestIDCtxKey struct{}

// XRequestIDKey is metadata key name for request ID
var XRequestIDKey = "x-request-id"

func UnaryServerInterceptor(logger log.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		reqId := handleRequestID(ctx)

		childCtx := lib.ContextWithValue(ctx, XRequestIDCtxKey{}, reqId)
		childCtx = lib.ContextWithValue(childCtx, log.LoggerCtxKey{}, logger.Log.With().Str(XRequestIDKey, reqId).Logger())

		return handler(childCtx, req)
	}
}

func handleRequestID(ctx context.Context) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return lib.NewUUID()
	}

	header, ok := md[XRequestIDKey]
	if !ok || len(header) == 0 {
		return lib.NewUUID()
	}

	requestID := header[0]
	if requestID == "" {
		return lib.NewUUID()
	}

	return requestID
}
