package store

import "github.com/bogdanrat/aws-serverless-poc/contracts/models"

type Store interface {
	GetAll() ([]*models.Book, error)
}
