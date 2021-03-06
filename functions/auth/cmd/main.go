package main

import (
	"context"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/bogdanrat/aws-serverless-poc/contracts/common"
	"github.com/bogdanrat/aws-serverless-poc/functions/auth/pkg/authenticator/cognitoauthenticator"
	"github.com/bogdanrat/aws-serverless-poc/functions/auth/pkg/handler"
	"os"
)

var (
	h *handler.Handler
)

func main() {
	lambda.Start(h.Handle)
}

func init() {
	// AWS Config
	cfg, _ := config.LoadDefaultConfig(context.TODO(), config.WithRegion(os.Getenv(common.AWSRegionEnvironmentVariable)))

	// Authenticator
	auth := cognitoauthenticator.New(cfg)

	// Handler
	h = handler.New(auth)
}
