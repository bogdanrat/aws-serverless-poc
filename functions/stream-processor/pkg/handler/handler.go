package handler

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/bogdanrat/aws-serverless-poc/functions/stream-processor/pkg/common"
	"github.com/bogdanrat/aws-serverless-poc/lib/store"
	"log"
)

type Handler struct{}

func New() *Handler {
	return &Handler{}
}

func (h *Handler) Handle(ctx context.Context, event events.DynamoDBEvent) {
	for _, record := range event.Records {
		switch record.EventName {
		case common.DynamoDBInsertEventName:
			h.handleInsertStreamRecord(record.Change)
		case common.DynamoDBModifyEventName:
			h.handleModifyStreamRecord(record.Change)
		case common.DynamoDBRemoveEventName:
			h.handleRemoveStreamRecord(record.Change)
		}
	}
}

func (h *Handler) handleInsertStreamRecord(record events.DynamoDBStreamRecord) {
	author := getBookAuthor(record)
	title := getBookTitle(record)

	log.Printf("New book published: %s - %s\n", author, title)
}

func (h *Handler) handleModifyStreamRecord(record events.DynamoDBStreamRecord) {
	author := getBookAuthor(record)
	title := getBookTitle(record)
	newInStockFormats, outOfStockFormats := getFormatsAvailability(record)

	if len(newInStockFormats) > 0 {
		log.Printf("New available formats for %s - %s:\n", author, title)
		for format := range newInStockFormats {
			log.Println(format)
		}
	}

	if len(outOfStockFormats) > 0 {
		log.Printf("Formats that went out of stock for %s - %s:\n", author, title)
		for format := range outOfStockFormats {
			log.Println(format)
		}
	}
}

func (h *Handler) handleRemoveStreamRecord(record events.DynamoDBStreamRecord) {
	author := getBookAuthor(record)
	title := getBookTitle(record)

	log.Printf("Book went out of stock: %s - %s\n", author, title)
}

func getFormatsAvailability(record events.DynamoDBStreamRecord) (map[string]bool, map[string]bool) {
	updatedFormats := make(map[string]bool)
	oldFormats := make(map[string]bool)
	oldFormatsAttributeValue := record.OldImage[store.FormatsTableAttributeName]
	if !oldFormatsAttributeValue.IsNull() {
		for formatType := range oldFormatsAttributeValue.Map() {
			oldFormats[formatType] = true
		}
	}

	newInStockFormats := make(map[string]bool)
	outOfStockFormats := make(map[string]bool)

	newFormatsAttributeValue := record.NewImage[store.FormatsTableAttributeName]
	if !newFormatsAttributeValue.IsNull() {
		for formatType := range newFormatsAttributeValue.Map() {
			updatedFormats[formatType] = true
			if !oldFormats[formatType] {
				newInStockFormats[formatType] = true
			}
		}
	}

	for format := range oldFormats {
		if !updatedFormats[format] {
			outOfStockFormats[format] = true
		}
	}

	return newInStockFormats, outOfStockFormats
}

func getBookAuthor(record events.DynamoDBStreamRecord) string {
	return record.Keys[store.AuthorTableAttributeName].String()
}
func getBookTitle(record events.DynamoDBStreamRecord) string {
	return record.Keys[store.TitleTableAttributeName].String()
}
