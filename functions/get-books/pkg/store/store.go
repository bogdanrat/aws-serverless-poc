package store

import "github.com/bogdanrat/aws-serverless-poc/functions/get-books/pkg/models"

type Store interface {
	GetAll() ([]*models.Book, error)
}
