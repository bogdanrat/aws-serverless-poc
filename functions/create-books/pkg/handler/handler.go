package handler

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/bogdanrat/aws-serverless-poc/contracts/common"
	"github.com/bogdanrat/aws-serverless-poc/contracts/models"
	"github.com/bogdanrat/aws-serverless-poc/lib/logger"
	"github.com/bogdanrat/aws-serverless-poc/lib/store"
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
		return h.apiResponse(http.StatusMethodNotAllowed, []byte(common.MethodNotAllowedErr.Error()))
	}

	books := make([]*models.Book, 0)
	err := json.Unmarshal([]byte(req.Body), &books)
	if err != nil {
		return h.apiResponse(http.StatusBadRequest, []byte(common.InvalidPayloadErr.Error()))
	}

	if err := h.Store.PutMany(books); err != nil {
		return h.apiResponse(http.StatusInternalServerError, []byte(fmt.Sprintf("%s: %v", common.DynamoDBActionErr, err)))
	}

	if err := h.logRequest(req, common.InsertedBooksMetricName); err != nil {
		return h.apiResponse(http.StatusInternalServerError, []byte(fmt.Sprintf("%s: %s", common.CWLogsErr, err)))
	}

	return h.apiResponse(http.StatusCreated, nil)
}

func (h *Handler) apiResponse(statusCode int, body []byte) (*events.APIGatewayProxyResponse, error) {
	response := &events.APIGatewayProxyResponse{
		StatusCode: statusCode,
		Headers: map[string]string{
			common.ContentTypeHeader: common.ContentTypeApplicationJSON,
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
