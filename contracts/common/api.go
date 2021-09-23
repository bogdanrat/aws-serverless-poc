package common

import "errors"

const (
	HttpGetMethod    = "get"
	HttpPostMethod   = "post"
	HttpPutMethod    = "put"
	HttpPatchMethod  = "patch"
	HttpDeleteMethod = "delete"
)

const (
	ContentTypeHeader          = "Content-Type"
	ContentTypeApplicationJSON = "application/json"
)

var (
	MethodNotAllowedErr = errors.New("method not allowed")
	InvalidPayloadErr   = errors.New("invalid payload")
	DynamoDBActionErr   = errors.New("dynamodb error")
	CWLogsErr           = errors.New("could not log request")
)
