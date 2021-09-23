package common

import "errors"

var (
	MethodNotAllowedErr  = errors.New("method not allowed")
	InvalidPayloadErr    = errors.New("invalid payload")
	DynamoDBActionErr    = errors.New("dynamodb error")
	DynamoDBUnmarshalErr = errors.New("error unmarshalling dynamodb items")
	CWLogsErr            = errors.New("could not log request")
	MarshalBooksErr      = errors.New("error marshalling books")
	SNSPublishErr        = errors.New("publish sns error")
)
