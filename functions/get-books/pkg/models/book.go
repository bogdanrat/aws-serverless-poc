package models

type Book struct {
	Title    string `json:"title" dynamodbav:"Title"`
	Author   string `json:"author" dynamodbav:"Author"`
	Category string `json:"category" dynamodbav:"Category"`
	Formats  string `json:"formats" dynamodbav:"Formats"`
}
