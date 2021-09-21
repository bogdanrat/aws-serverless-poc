package dynamostore

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/bogdanrat/aws-serverless-poc/contracts/common"
	"github.com/bogdanrat/aws-serverless-poc/contracts/models"
	"github.com/bogdanrat/aws-serverless-poc/functions/get-books/pkg/store"
	"os"
)

type DynamoStore struct {
	Client    *dynamodb.Client
	TableName string
}

func New(cfg aws.Config) store.Store {
	return &DynamoStore{
		Client:    dynamodb.NewFromConfig(cfg),
		TableName: os.Getenv(common.BooksTableNameEnvironmentVariable),
	}
}

func (s *DynamoStore) GetAll() ([]*models.Book, error) {
	scanInput := &dynamodb.ScanInput{
		TableName: aws.String(s.TableName),
	}

	output, err := s.Client.Scan(context.Background(), scanInput)
	if err != nil {
		return nil, fmt.Errorf("error scanning dynamodb table: %v", err)
	}

	books := make([]*models.Book, 0)

	err = attributevalue.UnmarshalListOfMaps(output.Items, &books)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling dynamodb items: %v", err)
	}

	return books, nil
}
