package common

import "errors"

var (
	UserAlreadyExistsErr = errors.New("user already exists")
	UserNotFoundErr      = errors.New("user not found")
	UserNotConfirmedErr  = errors.New("user not confirmed")
	NotAuthorizedErr     = errors.New("not authorized")
)
