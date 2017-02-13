package crawler

import (
	"fmt"

	"github.com/PuerkitoBio/goquery"
)

func getLineOperation(doc *goquery.Document) Operation {
	operationString := doc.Find(".schedule_active_list_active_tab").Text()
	operation, err := convertToOperation(operationString)
	if err != nil {
		panic(fmt.Errorf("Line MUST have of required operation types in order to be processed, given: %v", operationString))

	}
	return operation
}
