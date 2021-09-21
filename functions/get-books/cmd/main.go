package main

import (
	"context"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/bogdanrat/aws-serverless-poc/functions/get-books/pkg/handler"
	"github.com/bogdanrat/aws-serverless-poc/functions/get-books/pkg/logger/cwlogger"
	"github.com/bogdanrat/aws-serverless-poc/functions/get-books/pkg/store/dynamostore"
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
	cfg, _ := config.LoadDefaultConfig(context.TODO(), config.WithRegion(os.Getenv("REGION")))

	// DynamoDB
	store := dynamostore.New(cfg)

	// CW Metrics
	logger := cwlogger.New(cfg)

	// Handler
	h = handler.New(store, logger)
}
