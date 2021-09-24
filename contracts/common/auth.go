package common

import "errors"

const (
	AUTH_PARAMETER_USERNAME = "USERNAME"
	AUTH_PARAMETER_PASSWORD = "PASSWORD"
)

var (
	UserAlreadyExistsErr = errors.New("user already exists")
	UserNotFoundErr      = errors.New("user not found")
	UserNotConfirmedErr  = errors.New("user not confirmed")
	NotAuthorizedErr     = errors.New("not authorized")
)
