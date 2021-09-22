package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/bogdanrat/aws-serverless-poc/functions/stream-processor/pkg/handler"
)

var (
	h *handler.Handler
)

func main() {
	lambda.Start(h.Handle)
}

func init() {
	// Handler
	h = handler.New()
}
