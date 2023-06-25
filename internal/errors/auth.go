package errors

import "github.com/pkg/errors"

var (
	ErrAuthIncorrectCredentials = errors.New("login: incorrect credentials")
	ErrAuthInvalidCredentials   = errors.New("login: invalid credentials")
)

var (
	ErrAuthAccessTokenIncorrect = errors.New("session: incorrect access token")
	ErrAuthAccessTokenExpired   = errors.New("session: access token expired")
	ErrAuthPermissionDenied     = errors.New("session: permission denied")
)
