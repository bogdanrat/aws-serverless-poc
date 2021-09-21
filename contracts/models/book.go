package models

type Book struct {
	Title    string            `json:"title" dynamodbav:"Title"`
	Author   string            `json:"author" dynamodbav:"Author"`
	Category string            `json:"category,omitempty" dynamodbav:"Category"`
	Formats  map[string]string `json:"formats,omitempty" dynamodbav:"Formats"`
}
