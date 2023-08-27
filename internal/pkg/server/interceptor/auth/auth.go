package auth

import (
	"context"
	"github.com/kirychukyurii/wasker/internal/lib"
	"github.com/kirychukyurii/wasker/internal/pkg/server/interceptor/requestid"
	"strings"

	v1directorypb "github.com/kirychukyurii/wasker/gen/go/directory/v1"
	"github.com/kirychukyurii/wasker/internal/app/directory/controller"
	"github.com/kirychukyurii/wasker/internal/errors"
	"github.com/kirychukyurii/wasker/internal/pkg/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

var skipAuthServices = []string{
	v1directorypb.AuthService_ServiceDesc.ServiceName,
}

var (
	headerAuthorize = "authorization"
	typeAuthorize   = "bearer"
)

// UnaryServerInterceptor returns a server interceptor function to authenticate && authorize unary RPC
func UnaryServerInterceptor(logger log.Logger, controller controller.Controllers) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		if ok := skipAuthInterceptor(info.FullMethod); !ok {
			return handler(ctx, req)
		}

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return "", errors.ErrRequestMissingMetadata
		}

		token, err := authFromMD(ctx, md[headerAuthorize], typeAuthorize)
		if err != nil {
			return nil, err
		}

		userId, err := controller.Auth.Authn(ctx, token)
		if err != nil {
			return nil, err
		}

		ctx = context.WithValue(ctx, log.LoggerCtxKey{}, logger.FromContext(ctx).Log.With().Int64("user", userId).Logger())

		service, method := splitFullMethodName(info.FullMethod)
		ok, err = controller.Auth.Authz(ctx, userId, service, method)
		if err != nil || !ok {
			return nil, err
		}

		return handler(ctx, req)
	}
}

// authFromMD is a helper function for extracting the :authorization header from the gRPC metadata of the request.
//
// It expects the `:authorization` header to be of a certain scheme (e.g. `basic`, `bearer`), in a
// case-insensitive format (see rfc2617, sec 1.2). If no such authorization is found, or the token
// is of wrong scheme, an error with gRPC status `Unauthenticated` is returned.
func authFromMD(ctx context.Context, md []string, expectedScheme string) (string, error) {
	if len(md) < 1 {
		return "", errors.NewUnauthenticatedError(errors.AppError{
			Message: errors.ErrAuthAccessTokenIncorrect.Error(),
			Details: errors.AppErrorDetail{
				Err:       errors.ErrAuthAccessTokenIncorrect,
				ErrReason: "NULL_METADATA",
				ErrDomain: "interceptor.auth.from_metadata",
				RequestId: lib.FromContext(ctx, requestid.XRequestIDCtxKey{}).(string),
			},
		})
	}

	if md[0] == "" {
		return "", errors.NewUnauthenticatedError(errors.AppError{
			Message: errors.ErrAuthAccessTokenIncorrect.Error(),
			Details: errors.AppErrorDetail{
				Err:       errors.ErrAuthAccessTokenIncorrect,
				ErrReason: "EMPTY_METADATA",
				ErrDomain: "interceptor.auth.from_metadata",
				RequestId: lib.FromContext(ctx, requestid.XRequestIDCtxKey{}).(string),
			},
		})
	}

	scheme, token, found := strings.Cut(md[0], " ")
	if !found {
		return "", errors.NewUnauthenticatedError(errors.AppError{
			Message: errors.ErrAuthAccessTokenIncorrect.Error(),
			Details: errors.AppErrorDetail{
				Err:       errors.ErrAuthAccessTokenIncorrect,
				ErrReason: "NULL_AUTHORIZATION_HEADER",
				ErrDomain: "interceptor.auth.from_metadata",
				RequestId: lib.FromContext(ctx, requestid.XRequestIDCtxKey{}).(string),
			},
		})
	}

	if !strings.EqualFold(scheme, expectedScheme) {
		return "", errors.NewUnauthenticatedError(errors.AppError{
			Message: errors.ErrAuthAccessTokenIncorrect.Error(),
			Details: errors.AppErrorDetail{
				Err:       errors.ErrAuthAccessTokenIncorrect,
				ErrReason: "INVALID_TOKEN_TYPE",
				ErrDomain: "interceptor.auth.from_metadata",
				RequestId: lib.FromContext(ctx, requestid.XRequestIDCtxKey{}).(string),
			},
		})
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

func splitFullMethodName(fullMethod string) (string, string) {
	fullMethod = strings.TrimPrefix(fullMethod, "/") // remove leading slash
	if i := strings.Index(fullMethod, "/"); i >= 0 {
		return fullMethod[:i], fullMethod[i+1:]
	}

	return "unknown", "unknown"
}
