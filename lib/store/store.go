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
	PutMany(books []*models.Book) error
	Search(map[string]string) ([]*models.Book, error)
}
