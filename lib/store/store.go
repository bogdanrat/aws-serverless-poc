package store

import (
	"errors"
	"github.com/bogdanrat/aws-serverless-poc/contracts/models"
)

const (
	TitleTableAttributeName    = "Title"
	AuthorTableAttributeName   = "Author"
	CategoryTableAttributeName = "Category"
	FormatsTableAttributeName  = "Formats"
)

var (
	EmptyQueryErr = errors.New("empty query")
)

type Store interface {
	GetAll() ([]*models.Book, error)
	PutMany([]*models.Book) error
	Update(*models.Book, bool) (*models.Book, error)
	DeleteMany([]*models.Book, string) error
	Search(map[string]string) ([]*models.Book, error)
}
