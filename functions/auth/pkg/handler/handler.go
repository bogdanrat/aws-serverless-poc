package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/bogdanrat/aws-serverless-poc/contracts/common"
	"github.com/bogdanrat/aws-serverless-poc/contracts/models"
	"github.com/bogdanrat/aws-serverless-poc/functions/auth/pkg/authenticator"
	"net/http"
	"strings"
)

type Handler struct {
	Authenticator authenticator.Authenticator
}

func New(authenticator authenticator.Authenticator) *Handler {
	return &Handler{
		Authenticator: authenticator,
	}
}

func (h *Handler) Handle(req *events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	if strings.ToLower(req.HTTPMethod) != common.HttpPostMethod {
		return h.apiResponse(http.StatusMethodNotAllowed, []byte(common.MethodNotAllowedErr.Error()))
	}

	switch strings.ToLower(req.Path) {
	case common.SignUpPath:
		return h.handleSignUp(req)
	case common.LogInPath:
		return h.handleLogIn(req)
	default:
		return h.apiResponse(http.StatusNotFound, []byte(common.PathNotFoundErr.Error()))
	}
}

func (h *Handler) handleSignUp(req *events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	user := &models.User{}

	if err := json.Unmarshal([]byte(req.Body), user); err != nil {
		return h.apiResponse(http.StatusBadRequest, []byte(common.InvalidPayloadErr.Error()))
	}

	if err := h.Authenticator.SignUp(user); err != nil {
		if errors.Is(err, common.UserAlreadyExistsErr) {
			return h.apiResponse(http.StatusBadRequest, []byte(common.EmailAlreadyRegistered))
		}
		return h.apiResponse(http.StatusInternalServerError, []byte(err.Error()))
	}

	return h.apiResponse(http.StatusCreated, []byte(fmt.Sprintf("%s; %s", common.SignUpSuccessFul, common.CheckEmailForVerificationLink)))
}

func (h *Handler) handleLogIn(req *events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	user := &models.User{}

	if err := json.Unmarshal([]byte(req.Body), user); err != nil {
		return h.apiResponse(http.StatusBadRequest, []byte(common.InvalidPayloadErr.Error()))
	}

	token, err := h.Authenticator.LogIn(user)
	if err != nil {
		if errors.Is(err, common.UserNotFoundErr) {
			return h.apiResponse(http.StatusNotFound, []byte(common.EmailNotFound))
		} else if errors.Is(err, common.UserNotConfirmedErr) {
			return h.apiResponse(http.StatusForbidden, []byte(common.EmailNotConfirmed))
		} else if errors.Is(err, common.NotAuthorizedErr) {
			return h.apiResponse(http.StatusUnauthorized, []byte(common.WrongPassword))
		}
		return h.apiResponse(http.StatusInternalServerError, []byte(err.Error()))
	}

	response, err := json.Marshal(token)
	if err != nil {
		return h.apiResponse(http.StatusInternalServerError, []byte(fmt.Sprintf("%s: %s", common.TokenMarshalErr, err)))
	}

	return h.apiResponse(http.StatusOK, response)
}

func (h *Handler) apiResponse(statusCode int, body []byte) (*events.APIGatewayProxyResponse, error) {
	response := &events.APIGatewayProxyResponse{
		StatusCode: statusCode,
		Headers:    common.GetCorsHeaders(),
	}
	if body != nil {
		response.Body = string(body)
	}

	return response, nil
}
