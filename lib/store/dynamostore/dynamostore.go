package dynamostore

import (
	"context"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/bogdanrat/aws-serverless-poc/contracts/common"
	"github.com/bogdanrat/aws-serverless-poc/contracts/models"
	"github.com/bogdanrat/aws-serverless-poc/lib/store"
	"log"
	"os"
	"strings"
)

type DynamoStore struct {
	Client            *dynamodb.Client
	TableName         string
	CategoryIndexName string
	TitleIndexName    string
}

func New(cfg aws.Config) store.Store {
	return &DynamoStore{
		Client:            dynamodb.NewFromConfig(cfg),
		TableName:         os.Getenv(common.BooksTableNameEnvironmentVariable),
		CategoryIndexName: os.Getenv(common.BooksCategoryIndexNameEnvironmentVariable),
		TitleIndexName:    os.Getenv(common.BooksTitleIndexNameEnvironmentVariable),
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

func (s *DynamoStore) Search(queryParams map[string]string) ([]*models.Book, error) {
	queryInput, err := s.generateQueryInput(queryParams)
	if err != nil {
		return nil, err
	}

	output, err := s.Client.Query(context.Background(), queryInput)
	if err != nil {
		return nil, fmt.Errorf("error querying dynamodb table: %v", err)
	}

	books := make([]*models.Book, 0)

	err = attributevalue.UnmarshalListOfMaps(output.Items, &books)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling dynamodb items: %v", err)
	}

	return books, nil
}

func (s *DynamoStore) generateQueryInput(queryParams map[string]string) (*dynamodb.QueryInput, error) {
	queryInput := &dynamodb.QueryInput{
		TableName: aws.String(s.TableName),
	}

	category := queryParams[strings.ToLower(store.CategoryTableAttributeName)]
	author := queryParams[strings.ToLower(store.AuthorTableAttributeName)]
	title := queryParams[strings.ToLower(store.TitleTableAttributeName)]

	var keyConditionExpression string

	if author == "" {
		// query by category, filter by title
		if category != "" {
			queryInput.IndexName = aws.String(s.CategoryIndexName)
			keyConditionExpression = fmt.Sprintf("%s = :category", store.CategoryTableAttributeName)
			queryInput.ExpressionAttributeValues = map[string]types.AttributeValue{
				":category": &types.AttributeValueMemberS{Value: category},
			}

			if title != "" {
				queryInput.FilterExpression = aws.String(fmt.Sprintf("#t = :title"))
				queryInput.ExpressionAttributeNames["#t"] = store.TitleTableAttributeName
				queryInput.ExpressionAttributeValues[":title"] = &types.AttributeValueMemberS{Value: title}
			}
		} else if title != "" {
			// query by title
			queryInput.IndexName = aws.String(s.TitleIndexName)
			keyConditionExpression = fmt.Sprintf("%s = :title", store.TitleTableAttributeName)
			queryInput.ExpressionAttributeValues = map[string]types.AttributeValue{
				":title": &types.AttributeValueMemberS{Value: title},
			}
		}
	} else {
		keyConditionExpression = fmt.Sprintf(`%s = :author`, store.AuthorTableAttributeName)
		queryInput.ExpressionAttributeValues = map[string]types.AttributeValue{
			":author": &types.AttributeValueMemberS{Value: author},
		}

		if title != "" {
			keyConditionExpression = fmt.Sprintf("%s AND %s = :title", keyConditionExpression, store.TitleTableAttributeName)
			queryInput.ExpressionAttributeValues[":title"] = &types.AttributeValueMemberS{Value: title}
		}
	}

	if keyConditionExpression == "" {
		return nil, errors.New("invalid query params")
	}

	queryInput.KeyConditionExpression = aws.String(keyConditionExpression)

	log.Printf("keyConditionExpression: %s\n", keyConditionExpression)
	log.Printf("ExpressionAttributeNames: %v\n", queryInput.ExpressionAttributeNames)
	log.Printf("ExpressionAttributeValues: %v\n", queryInput.ExpressionAttributeValues)

	return queryInput, nil
}
