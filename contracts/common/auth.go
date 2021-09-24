package common

import "errors"

const (
	AUTH_PARAMETER_USERNAME       = "USERNAME"
	AUTH_PARAMETER_PASSWORD       = "PASSWORD"
	EmailNotFound                 = "email not found"
	EmailNotConfirmed             = "email not confirmed"
	SignUpSuccessFul              = "sign up successful"
	CheckEmailForVerificationLink = "please check your email for verification link"
	WrongPassword                 = "wrong password"
)

var (
	UserAlreadyExistsErr = errors.New("user already exists")
	UserNotFoundErr      = errors.New("user not found")
	UserNotConfirmedErr  = errors.New("user not confirmed")
	NotAuthorizedErr     = errors.New("not authorized")
)
