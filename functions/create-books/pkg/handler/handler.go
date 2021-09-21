package handler

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/bogdanrat/aws-serverless-poc/contracts/common"
	"github.com/bogdanrat/aws-serverless-poc/contracts/models"
	"github.com/bogdanrat/aws-serverless-poc/functions/create-books/pkg/store"
	"github.com/bogdanrat/aws-serverless-poc/lib/logger"
	"net/http"
	"strings"
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
	if !strings.EqualFold(req.HTTPMethod, common.HttpPostMethod) {
		return h.apiResponse(http.StatusMethodNotAllowed, []byte("method not allowed"))
	}

	books := make([]*models.Book, 0)
	err := json.Unmarshal([]byte(req.Body), &books)
	if err != nil {
		return h.apiResponse(http.StatusBadRequest, []byte(fmt.Sprintf("invalid payload")))
	}

	if err := h.Store.PutMany(books); err != nil {
		return h.apiResponse(http.StatusInternalServerError, []byte(fmt.Sprintf("error inserting into dynamodb: %v", err)))
	}

	if err := h.logRequest(req, common.InsertedBooksMetricName); err != nil {
		return h.apiResponse(http.StatusInternalServerError, []byte(fmt.Sprintf("could not log request: %s", err)))
	}

	return h.apiResponse(http.StatusCreated, nil)
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
