package errors

import (
	"encoding/json"

	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
)

type AppError struct {
	Code    codes.Code     `json:"code"`
	Message string         `json:"message"`
	Details AppErrorDetail `json:"details"`
}

type AppErrorDetail struct {
	Err       error  `json:"error"`
	ErrDomain string `json:"domain"`
	ErrReason string `json:"reason"`
	RequestId string `json:"request_id"`
}

// Error Returns Message if Details.Err is nil.
func (e AppError) Error() string {
	if err := e.Details.Err; err != nil {
		return err.Error()
	}

	return e.Message
}

func (e AppError) Msg() string {
	return e.Message
}

func (e AppError) ToJson() string {
	b, _ := json.Marshal(e)

	return string(b)
}

func NewInternalError(appError AppError) *AppError {
	return &AppError{
		Code:    codes.Internal,
		Message: appError.Message,
		Details: AppErrorDetail{
			Err:       appError.Details.Err,
			ErrDomain: appError.Details.ErrDomain,
			ErrReason: appError.Details.ErrReason,
			RequestId: appError.Details.RequestId,
		},
	}
}

func NewUnauthenticatedError(appError AppError) *AppError {
	return &AppError{
		Code:    codes.Unauthenticated,
		Message: appError.Message,
		Details: AppErrorDetail{
			Err:       appError.Details.Err,
			ErrDomain: appError.Details.ErrDomain,
			ErrReason: appError.Details.ErrReason,
			RequestId: appError.Details.RequestId,
		},
	}
}

func NewForbiddenError(appError AppError) *AppError {
	return &AppError{
		Code:    codes.PermissionDenied,
		Message: appError.Message,
		Details: AppErrorDetail{
			Err:       appError.Details.Err,
			ErrDomain: appError.Details.ErrDomain,
			ErrReason: appError.Details.ErrReason,
			RequestId: appError.Details.RequestId,
		},
	}
}

func NewNotFoundError(appError AppError) *AppError {
	return &AppError{
		Code:    codes.NotFound,
		Message: appError.Message,
		Details: AppErrorDetail{
			Err:       appError.Details.Err,
			ErrDomain: appError.Details.ErrDomain,
			ErrReason: appError.Details.ErrReason,
			RequestId: appError.Details.RequestId,
		},
	}
}

func NewBadRequestError(appError AppError) *AppError {
	return &AppError{
		Code:    codes.Internal,
		Message: appError.Message,
		Details: AppErrorDetail{
			Err:       appError.Details.Err,
			ErrDomain: appError.Details.ErrDomain,
			ErrReason: appError.Details.ErrReason,
			RequestId: appError.Details.RequestId,
		},
	}
}

// Alias
// noinspection GoUnusedGlobalVariable
var (
	Is           = errors.Is
	As           = errors.As
	New          = errors.New
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
