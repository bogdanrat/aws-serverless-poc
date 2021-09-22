package store

import "github.com/bogdanrat/aws-serverless-poc/contracts/models"

const (
	TitleTableAttributeName    = "Title"
	AuthorTableAttributeName   = "Author"
	CategoryTableAttributeName = "Category"
	FormatsTableAttributeName  = "Formats"
)

type Store interface {
	GetAll() ([]*models.Book, error)
	PutMany(books []*models.Book) error
	Search(map[string]string) ([]*models.Book, error)
}
