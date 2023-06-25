package errors

import (
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Alias
// noinspection GoUnusedGlobalVariable
var (
	Is = errors.Is
	As = errors.As
	//New          = errors.New
	Errorf       = errors.Errorf
	Unwrap       = errors.Unwrap
	Wrap         = errors.Wrap
	Wrapf        = errors.Wrapf
	WithStack    = errors.WithStack
	WithMessage  = errors.WithMessage
	WithMessagef = errors.WithMessagef
)

var (
	ErrRequestMissingMetadata = errors.New("request: missing metadata")
)

func New(err error) error {
	switch {
	case Is(err, ErrAuthPermissionDenied):
		return status.Error(codes.PermissionDenied, err.Error())
	case Is(err, ErrAuthInvalidCredentials) || Is(err, ErrAuthIncorrectCredentials) ||
		Is(err, ErrAuthAccessTokenExpired) || Is(err, ErrAuthAccessTokenIncorrect):
		return status.Error(codes.Unauthenticated, err.Error())
	case Is(err, ErrDatabaseRecordNotFound):
		return status.Error(codes.NotFound, err.Error())
	case Is(err, ErrDatabaseInternalError):
		return status.Error(codes.Internal, err.Error())
	case Is(err, ErrRequestMissingMetadata):
		return status.Error(codes.InvalidArgument, ErrRequestMissingMetadata.Error())
	}

	return err
}
