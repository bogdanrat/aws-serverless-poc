package authenticator

import "github.com/bogdanrat/aws-serverless-poc/contracts/models"

type Authenticator interface {
	SignUp(*models.User) error
	LogIn(*models.User) (*models.Token, error)
}
