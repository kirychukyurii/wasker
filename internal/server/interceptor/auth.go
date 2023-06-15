package interceptor

import (
	"context"
	"strings"

	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	v1directorypb "github.com/kirychukyurii/wasker/gen/go/directory/v1"
	"github.com/kirychukyurii/wasker/internal/controller"
	"github.com/kirychukyurii/wasker/internal/model"
	"github.com/kirychukyurii/wasker/internal/pkg/log"
)

var skipAuthServices = []string{
	v1directorypb.AuthService_ServiceDesc.ServiceName,
}

var (
	headerAuthorize = "authorization"
	typeAuthorize   = "bearer"
)

func AuthUnaryServerInterceptor(logger log.Logger, controller controller.Controllers) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		service, method := splitFullMethodName(info.FullMethod)
		logger.Log.Debug().Str("service", service).Str("method", method).Msg("test")
		if ok := skipAuthInterceptor(service); !ok {
			return handler(ctx, req)
		}

		token, err := authFromMD(ctx, typeAuthorize)
		if err != nil {
			logger.Log.Error().Err(err).Msg("extract authorization header")
			return nil, err
		}

		err = controller.Auth.CheckToken(ctx, token, serviceToScope(service), methodToPermission(method))
		if err != nil {
			logger.Log.Error().Err(errors.Wrap(err, "controller.V1alpha1.Auth.CheckToken()")).Msg("check token")
			return nil, status.Error(codes.Unauthenticated, errors.Wrap(err, "controller.V1alpha1.Auth.CheckToken()").Error())
		}

		return handler(ctx, req)
	}
}

// AuthFromMD is a helper function for extracting the :authorization header from the gRPC metadata of the request.
//
// It expects the `:authorization` header to be of a certain scheme (e.g. `basic`, `bearer`), in a
// case-insensitive format (see rfc2617, sec 1.2). If no such authorization is found, or the token
// is of wrong scheme, an error with gRPC status `Unauthenticated` is returned.
func authFromMD(ctx context.Context, expectedScheme string) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", status.Errorf(codes.InvalidArgument, "missing metadata")
	}

	authorization := md["authorization"]
	if len(authorization) < 1 {
		return "", errors.New("received empty authorization token from client")
	}

	val := md.Get(headerAuthorize)
	if val[0] == "" {
		return "", status.Error(codes.Unauthenticated, "Request unauthenticated with "+expectedScheme)
	}

	scheme, token, found := strings.Cut(val[0], " ")
	if !found {
		return "", status.Error(codes.Unauthenticated, "Bad authorization string")
	}

	if !strings.EqualFold(scheme, expectedScheme) {
		return "", status.Error(codes.Unauthenticated, "Request unauthenticated with "+expectedScheme)
	}

	return token, nil
}

func splitFullMethodName(fullMethod string) (string, string) {
	fullMethod = strings.TrimPrefix(fullMethod, "/") // remove leading slash
	if i := strings.Index(fullMethod, "/"); i >= 0 {
		return fullMethod[:i], fullMethod[i+1:]
	}

	return "unknown", "unknown"
}

// skipAuthInterceptor setup auth matcher.
func skipAuthInterceptor(service string) bool {
	for _, s := range skipAuthServices {
		return s != service
	}

	return true
}

func serviceToScope(service string) (scope string) {
	switch service {
	case v1directorypb.UserService_ServiceDesc.ServiceName:
		scope = model.UserScope
	case v1directorypb.RoleService_ServiceDesc.ServiceName:
		scope = model.RoleScope
	}
	return scope
}

func methodToPermission(method string) (permission string) {
	if strings.Contains(method, model.CreatePermission) {
		permission = model.CreatePermission
	} else if strings.Contains(method, model.ReadPermission) {
		permission = model.ReadPermission
	} else if strings.Contains(method, model.UpdatePermission) {
		permission = model.UpdatePermission
	} else if strings.Contains(method, model.DeletePermission) {
		permission = model.DeletePermission
	}

	return permission
}
