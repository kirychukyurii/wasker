package interceptor

import (
	"context"
	"github.com/kirychukyurii/wasker/internal/config"
	"github.com/kirychukyurii/wasker/internal/pkg/log"
	"github.com/kirychukyurii/wasker/internal/server/interceptor/requestid"
	"google.golang.org/grpc/peer"
	"strings"
	"time"

	"google.golang.org/grpc"
)

var (
	ClientIPCtxKey = "client-ip"
	ServiceCtxKey  = "service"
	MethodCtxKey   = "method"
)

func ContextUnaryServerInterceptor(cfg config.Config, logger log.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		ctx, cancel := context.WithTimeout(ctx, time.Duration(cfg.Grpc.Timeout)*time.Second)
		defer cancel()

		ctx = withClientIP(ctx)
		ctx = withServiceInfo(ctx, info.FullMethod)
		ctx = withLogger(ctx, logger)

		return handler(ctx, req)
	}
}

func FromContext(ctx context.Context, key string) string {
	return ctx.Value(key).(string)
}

func withClientIP(ctx context.Context) context.Context {
	var clientIP string

	p, ok := peer.FromContext(ctx)
	if !ok {
		clientIP = "0.0.0.0"
	}

	clientAddr := p.Addr.String()
	if i := strings.Index(clientAddr, ":"); i >= 0 {
		clientIP = clientAddr[:i]
	}

	ctx = context.WithValue(ctx, ClientIPCtxKey, clientIP)

	return ctx
}

func withServiceInfo(ctx context.Context, fullMethod string) context.Context {
	service, method := splitFullMethodName(fullMethod)

	ctx = context.WithValue(ctx, ServiceCtxKey, service)
	ctx = context.WithValue(ctx, MethodCtxKey, method)

	return ctx
}

func splitFullMethodName(fullMethod string) (string, string) {
	fullMethod = strings.TrimPrefix(fullMethod, "/") // remove leading slash
	if i := strings.Index(fullMethod, "/"); i >= 0 {
		return fullMethod[:i], fullMethod[i+1:]
	}

	return "unknown", "unknown"
}

func withLogger(ctx context.Context, logger log.Logger) context.Context {
	service := FromContext(ctx, ServiceCtxKey)
	method := FromContext(ctx, MethodCtxKey)
	clientIP := FromContext(ctx, ClientIPCtxKey)

	subLogger := logger.Log.With().Str("service", service).Str("method", method).
		Str(requestid.DefaultXRequestIDKey, requestid.FromContext(ctx)).Str(ClientIPCtxKey, clientIP).Logger()
	ctx = context.WithValue(ctx, log.LoggerCtxKey{}, subLogger)

	return ctx
}
