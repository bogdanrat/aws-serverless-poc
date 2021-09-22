package dynamostore

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/bogdanrat/aws-serverless-poc/contracts/common"
	"github.com/bogdanrat/aws-serverless-poc/contracts/models"
	"github.com/bogdanrat/aws-serverless-poc/lib/store"
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

func (s *DynamoStore) PutMany(books []*models.Book) error {
	input := &dynamodb.BatchWriteItemInput{
		RequestItems: make(map[string][]types.WriteRequest),
	}

	for _, book := range books {
		writeRequest := types.WriteRequest{
			PutRequest: &types.PutRequest{
				Item: map[string]types.AttributeValue{
					store.TitleTableAttributeName:    &types.AttributeValueMemberS{Value: book.Title},
					store.AuthorTableAttributeName:   &types.AttributeValueMemberS{Value: book.Author},
					store.CategoryTableAttributeName: &types.AttributeValueMemberS{Value: book.Category},
				},
			},
		}

		formatsMap := make(map[string]types.AttributeValue)
		for formatKey, formatValue := range book.Formats {
			formatsMap[formatKey] = &types.AttributeValueMemberS{Value: formatValue}
		}
		writeRequest.PutRequest.Item[store.FormatsTableAttributeName] = &types.AttributeValueMemberM{Value: formatsMap}

		input.RequestItems[s.TableName] = append(input.RequestItems[s.TableName], writeRequest)
	}

	_, err := s.Client.BatchWriteItem(context.TODO(), input)
	if err != nil {
		return err
	}

	return nil
}