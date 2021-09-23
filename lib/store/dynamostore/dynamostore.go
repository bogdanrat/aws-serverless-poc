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
	"os"
	"strings"
)

const (
	authorExpressionAttributeName    = "#author"
	authorExpressionAttributeValue   = ":author"
	titleExpressionAttributeName     = "#title"
	titleExpressionAttributeValue    = ":title"
	categoryExpressionAttributeName  = "#category"
	categoryExpressionAttributeValue = ":category"
	formatsExpressionAttributeName   = "#formats"
	formatsExpressionAttributeValue  = ":formats"
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
		// if no search params were provided, return all books
		if errors.Is(err, store.EmptyQueryErr) {
			return s.GetAll()
		}
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
			keyConditionExpression = fmt.Sprintf("%s = %s", store.CategoryTableAttributeName, categoryExpressionAttributeValue)
			queryInput.ExpressionAttributeValues = map[string]types.AttributeValue{
				categoryExpressionAttributeValue: &types.AttributeValueMemberS{Value: category},
			}

			if title != "" {
				queryInput.FilterExpression = aws.String(fmt.Sprintf("%s = %s", titleExpressionAttributeName, titleExpressionAttributeValue))
				queryInput.ExpressionAttributeNames = map[string]string{
					titleExpressionAttributeName: store.TitleTableAttributeName,
				}
				queryInput.ExpressionAttributeValues[titleExpressionAttributeValue] = &types.AttributeValueMemberS{Value: title}
			}
		} else if title != "" {
			// query by title
			queryInput.IndexName = aws.String(s.TitleIndexName)
			keyConditionExpression = fmt.Sprintf("%s = %s", store.TitleTableAttributeName, titleExpressionAttributeValue)
			queryInput.ExpressionAttributeValues = map[string]types.AttributeValue{
				titleExpressionAttributeValue: &types.AttributeValueMemberS{Value: title},
			}
		}
	} else {
		// query by author (and by title) and filter by category
		keyConditionExpression = fmt.Sprintf(`%s = %s`, store.AuthorTableAttributeName, authorExpressionAttributeValue)
		queryInput.ExpressionAttributeValues = map[string]types.AttributeValue{
			authorExpressionAttributeValue: &types.AttributeValueMemberS{Value: author},
		}

		if title != "" {
			keyConditionExpression = fmt.Sprintf("%s AND %s = %s", keyConditionExpression, store.TitleTableAttributeName, titleExpressionAttributeValue)
			queryInput.ExpressionAttributeValues[titleExpressionAttributeValue] = &types.AttributeValueMemberS{Value: title}
		} else if category != "" {
			queryInput.FilterExpression = aws.String(fmt.Sprintf("%s = %s", categoryExpressionAttributeName, categoryExpressionAttributeValue))
			queryInput.ExpressionAttributeNames = map[string]string{
				categoryExpressionAttributeName: store.CategoryTableAttributeName,
			}
			queryInput.ExpressionAttributeValues[categoryExpressionAttributeValue] = &types.AttributeValueMemberS{Value: category}
		}
	}

	if keyConditionExpression == "" {
		return nil, store.EmptyQueryErr
	}

	queryInput.KeyConditionExpression = aws.String(keyConditionExpression)
	return queryInput, nil
}

func (s *DynamoStore) Update(book *models.Book, partial bool) (*models.Book, error) {
	updateInput, err := s.generateUpdateInput(book, partial)
	if err != nil {
		return nil, err
	}

	output, err := s.Client.UpdateItem(context.Background(), updateInput)
	if err != nil {
		return nil, err
	}

	updatedBook := &models.Book{}
	err = attributevalue.UnmarshalMap(output.Attributes, &updatedBook)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling dynamodb attributes map: %v", err)
	}

	return updatedBook, nil
}

func (s *DynamoStore) generateUpdateInput(book *models.Book, partial bool) (*dynamodb.UpdateItemInput, error) {
	formatsMap := make(map[string]types.AttributeValue)
	for formatKey, formatValue := range book.Formats {
		formatsMap[formatKey] = &types.AttributeValueMemberS{Value: formatValue}
	}

	updateInput := &dynamodb.UpdateItemInput{
		TableName: aws.String(s.TableName),
		Key: map[string]types.AttributeValue{
			store.AuthorTableAttributeName: &types.AttributeValueMemberS{Value: book.Author},
			store.TitleTableAttributeName:  &types.AttributeValueMemberS{Value: book.Title},
		},
		ReturnValues: types.ReturnValueAllNew,
	}

	var updateExpression string
	// full update
	if !partial {
		// update formats no matter what, even with blank values
		updateInput.ExpressionAttributeNames = map[string]string{
			formatsExpressionAttributeName: store.FormatsTableAttributeName,
		}
		updateInput.ExpressionAttributeValues = map[string]types.AttributeValue{
			formatsExpressionAttributeValue: &types.AttributeValueMemberM{Value: formatsMap},
		}

		updateExpression = fmt.Sprintf("SET %s = %s", formatsExpressionAttributeName, formatsExpressionAttributeValue)

		// The AttributeValue for a key attribute (Category GSI) cannot contain an empty string value
		// only update category if value was provided
		if book.Category != "" {
			updateInput.ExpressionAttributeNames[categoryExpressionAttributeName] = store.CategoryTableAttributeName
			updateInput.ExpressionAttributeValues[categoryExpressionAttributeValue] = &types.AttributeValueMemberS{Value: book.Category}
			updateExpression = fmt.Sprintf("%s, %s = %s", updateExpression, categoryExpressionAttributeName, categoryExpressionAttributeValue)
		}
	} else {
		// partial update
		// only update category if value was provided
		if book.Category != "" {
			updateInput.ExpressionAttributeNames = map[string]string{
				categoryExpressionAttributeName: store.CategoryTableAttributeName,
			}
			updateInput.ExpressionAttributeValues = map[string]types.AttributeValue{
				categoryExpressionAttributeValue: &types.AttributeValueMemberS{Value: book.Category},
			}
			updateExpression = fmt.Sprintf("SET %s = %s", categoryExpressionAttributeName, categoryExpressionAttributeValue)

			// only update formats if value was provided
			if len(formatsMap) > 0 {
				updateInput.ExpressionAttributeNames[formatsExpressionAttributeName] = store.FormatsTableAttributeName
				updateInput.ExpressionAttributeValues[formatsExpressionAttributeValue] = &types.AttributeValueMemberM{Value: formatsMap}
				updateExpression = fmt.Sprintf("SET %s = %s, %s = %s", categoryExpressionAttributeName, categoryExpressionAttributeValue, formatsExpressionAttributeName, formatsExpressionAttributeValue)
			}
		} else if len(formatsMap) > 0 {
			updateInput.ExpressionAttributeNames = map[string]string{
				formatsExpressionAttributeName: store.FormatsTableAttributeName,
			}
			updateInput.ExpressionAttributeValues = map[string]types.AttributeValue{
				formatsExpressionAttributeValue: &types.AttributeValueMemberM{Value: formatsMap},
			}
			updateExpression = fmt.Sprintf("SET %s = %s", formatsExpressionAttributeName, formatsExpressionAttributeValue)
		}
	}

	updateInput.UpdateExpression = aws.String(updateExpression)
	return updateInput, nil
}

func (s *DynamoStore) DeleteMany(books []*models.Book, requestToken string) error {
	transactInput := &dynamodb.TransactWriteItemsInput{
		TransactItems: make([]types.TransactWriteItem, 0),
		// If the original TransactWriteItems call was successful,
		// the subsequent TransactWriteItems calls with the same client token return successfully without making any changes.
		ClientRequestToken: aws.String(requestToken),
	}
	for _, book := range books {
		transactInput.TransactItems = append(transactInput.TransactItems, types.TransactWriteItem{
			Delete: &types.Delete{
				TableName: aws.String(s.TableName),
				Key: map[string]types.AttributeValue{
					store.AuthorTableAttributeName: &types.AttributeValueMemberS{Value: book.Author},
					store.TitleTableAttributeName:  &types.AttributeValueMemberS{Value: book.Title},
				},
			},
		})
	}

	_, err := s.Client.TransactWriteItems(context.TODO(), transactInput)
	if err != nil {
		return err
	}

	return nil
}
