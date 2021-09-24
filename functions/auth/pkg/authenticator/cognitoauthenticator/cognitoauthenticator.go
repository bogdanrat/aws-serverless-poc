package cognitoauthenticator

import (
	"context"
	"errors"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
	"github.com/bogdanrat/aws-serverless-poc/contracts/common"
	"github.com/bogdanrat/aws-serverless-poc/contracts/models"
	"github.com/bogdanrat/aws-serverless-poc/functions/auth/pkg/authenticator"
	"log"
	"os"
)

type CognitoAuthenticator struct {
	Client      *cognitoidentityprovider.Client
	UserPoolID  string
	AppClientID string
}

func New(cfg aws.Config) authenticator.Authenticator {
	return &CognitoAuthenticator{
		Client:      cognitoidentityprovider.NewFromConfig(cfg),
		UserPoolID:  os.Getenv(common.CognitoUserPoolIDEnvironmentVariable),
		AppClientID: os.Getenv(common.CognitoAppClientIDEnvironmentVariable),
	}
}

func (auth *CognitoAuthenticator) SignUp(user *models.User) error {
	signUpInput := &cognitoidentityprovider.SignUpInput{
		Username: aws.String(user.Email),
		Password: aws.String(user.Password),
		ClientId: aws.String(auth.AppClientID),
	}

	_, err := auth.Client.SignUp(context.TODO(), signUpInput)

	if err != nil {
		return unWrapError(err)
	}
	return nil
}

func (auth *CognitoAuthenticator) LogIn(user *models.User) (*models.Token, error) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("LogIn() recovered: %s\n", r)
		}
	}()

	authInput := &cognitoidentityprovider.InitiateAuthInput{
		AuthFlow: types.AuthFlowTypeUserPasswordAuth,
		AuthParameters: map[string]string{
			common.AUTH_PARAMETER_USERNAME: user.Email,
			common.AUTH_PARAMETER_PASSWORD: user.Password,
		},
		ClientId: aws.String(auth.AppClientID),
	}

	output, err := auth.Client.InitiateAuth(context.TODO(), authInput)
	if err != nil {
		return nil, unWrapError(err)
	}

	log.Printf("Authenticated: %s\n", *output.AuthenticationResult.IdToken)

	return &models.Token{
		ExpiresIn:    output.AuthenticationResult.ExpiresIn,
		AccessToken:  *output.AuthenticationResult.IdToken,
		RefreshToken: *output.AuthenticationResult.RefreshToken,
	}, nil
}

func unWrapError(err error) error {
	var userNotFound *types.UserNotFoundException
	var userExists *types.UsernameExistsException
	var userNotConfirmed *types.UserNotConfirmedException
	var notAuthorized *types.NotAuthorizedException

	if errors.As(err, &userNotFound) {
		return common.UserNotFoundErr
	} else if errors.As(err, &userExists) {
		return common.UserAlreadyExistsErr
	} else if errors.As(err, &userNotConfirmed) {
		return common.UserNotConfirmedErr
	} else if errors.As(err, &notAuthorized) {
		return common.NotAuthorizedErr
	}

	return err
}
