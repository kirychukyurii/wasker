package auth

import (
	"context"
	"fmt"
	"github.com/kirychukyurii/wasker/internal/server/interceptor"
	"strconv"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	v1directorypb "github.com/kirychukyurii/wasker/gen/go/directory/v1"
	"github.com/kirychukyurii/wasker/internal/controller"
	"github.com/kirychukyurii/wasker/internal/errors"
	"github.com/kirychukyurii/wasker/internal/pkg/log"
)

var skipAuthServices = []string{
	v1directorypb.AuthService_ServiceDesc.ServiceName,
}

var (
	headerAuthorize = "authorization"
	typeAuthorize   = "bearer"
)

func UnaryServerInterceptor(logger log.Logger, controller controller.Controllers) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		service := interceptor.FromContext(ctx, interceptor.ServiceCtxKey)
		method := interceptor.FromContext(ctx, interceptor.MethodCtxKey)

		if ok := skipAuthInterceptor(service); !ok {
			return handler(ctx, req)
		}

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return "", errors.New(errors.ErrRequestMissingMetadata)
		}

		token, err := authFromMD(md[headerAuthorize], typeAuthorize)
		if err != nil {
			return nil, status.Error(codes.Unauthenticated, err.Error())
		}

		userId, err := controller.Auth.VerifyToken(ctx, token)
		if err != nil {
			logger.FromContext(ctx).Log.Error().Err(err).Msg("verify token")
			return nil, status.Error(codes.Unauthenticated, "unauthenticated")
		}

		subLogger := logger.FromContext(ctx).Log.With().Str("user", strconv.FormatUint(userId, 10)).Logger()
		ctx = context.WithValue(ctx, log.LoggerCtxKey{}, subLogger)

		ok, err = controller.Auth.VerifyPermission(ctx, userId, service, method)
		if err != nil || !ok {
			logger.FromContext(ctx).Log.Error().Err(err).Msg("permission denied")

			return nil, status.Error(codes.PermissionDenied, fmt.Sprintf("denied: user_id=%d, service=%s, method=%s", userId, service, method))
		}

		return handler(ctx, req)
	}
}

// AuthFromMD is a helper function for extracting the :authorization header from the gRPC metadata of the request.
//
// It expects the `:authorization` header to be of a certain scheme (e.g. `basic`, `bearer`), in a
// case-insensitive format (see rfc2617, sec 1.2). If no such authorization is found, or the token
// is of wrong scheme, an error with gRPC status `Unauthenticated` is returned.
func authFromMD(md []string, expectedScheme string) (string, error) {
	if len(md) < 1 {
		return "", errors.New(errors.ErrAuthAccessTokenIncorrect)
	}

	if md[0] == "" {
		return "", errors.New(errors.ErrAuthAccessTokenIncorrect)
	}

	scheme, token, found := strings.Cut(md[0], " ")
	if !found {
		return "", errors.New(errors.ErrAuthAccessTokenIncorrect)
	}

	if !strings.EqualFold(scheme, expectedScheme) {
		return "", errors.New(errors.ErrAuthAccessTokenIncorrect)
	}

	return token, nil
}

// skipAuthInterceptor setup auth matcher.
func skipAuthInterceptor(service string) bool {
	for _, s := range skipAuthServices {
		return s != service
	}

	return true
}
