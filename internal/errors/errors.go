package errors

import (
	"github.com/pkg/errors"
)

type AppError struct {
	Code        int
	Id          string
	Status      string
	Description error
}

// Alias
// noinspection GoUnusedGlobalVariable
var (
	Is           = errors.Is
	As           = errors.As
	New          = errors.New
	Unwrap       = errors.Unwrap
	Wrap         = errors.Wrap
	Wrapf        = errors.Wrapf
	WithStack    = errors.WithStack
	WithMessage  = errors.WithMessage
	WithMessagef = errors.WithMessagef
)

// Database
// noinspection GoUnusedGlobalVariable
var (
	DatabaseInternalError  = errors.New("database internal error")
	DatabaseRecordNotFound = errors.New("database record not found")
)

// Auth
// noinspection GoUnusedGlobalVariable
var (
	AuthTokenInvalid      = errors.New("auth token is invalid")
	AuthTokenExpired      = errors.New("auth token is expired")
	AuthTokenNotValidYet  = errors.New("auth token not active yet")
	AuthTokenMalformed    = errors.New("auth token is malformed")
	AuthTokenGenerateFail = errors.New("failed to generate auth token")
)
