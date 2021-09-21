package store

import "github.com/bogdanrat/aws-serverless-poc/contracts/models"

const (
	TitleTableAttributeName    = "Title"
	AuthorTableAttributeName   = "Author"
	CategoryTableAttributeName = "Category"
	FormatsTableAttributeName  = "Formats"
)

type Store interface {
	PutMany(books []*models.Book) error
}
