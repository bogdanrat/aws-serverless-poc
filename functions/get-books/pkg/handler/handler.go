package handler

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/bogdanrat/aws-serverless-poc/functions/get-books/pkg/common"
	"github.com/bogdanrat/aws-serverless-poc/functions/get-books/pkg/logger"
	"github.com/bogdanrat/aws-serverless-poc/functions/get-books/pkg/store"
	"net/http"
)

type Handler struct {
	Store  store.Store
	Logger logger.MetricLogger
}

func New(store store.Store, logger logger.MetricLogger) *Handler {
	return &Handler{
		Store:  store,
		Logger: logger,
	}
}

func (h *Handler) Handle(req *events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	books, err := h.Store.GetAll()
	if err != nil {
		return h.apiResponse(http.StatusInternalServerError, []byte(err.Error()))
	}

	response, err := json.Marshal(books)
	if err != nil {
		return h.apiResponse(http.StatusInternalServerError, []byte(fmt.Sprintf("error marshalling books: %v", err)))
	}

	if err := h.logRequest(req, common.FetchedBooksMetricName); err != nil {
		return h.apiResponse(http.StatusInternalServerError, []byte(fmt.Sprintf("could not log request: %s", err)))
	}

	return h.apiResponse(http.StatusOK, response)
}

func (h *Handler) apiResponse(statusCode int, body []byte) (*events.APIGatewayProxyResponse, error) {
	response := &events.APIGatewayProxyResponse{
		StatusCode: statusCode,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}
	if body != nil {
		response.Body = string(body)
	}

	return response, nil
}

func (h *Handler) logRequest(req *events.APIGatewayProxyRequest, metricName string) error {
	alias, ok := req.Headers[common.LambdaAliasHeader]
	if !ok {
		alias = common.LambdaDefaultAlias
	}

	if err := h.Logger.PutMetric(alias, metricName); err != nil {
		return err
	}

	return nil
}
