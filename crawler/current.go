package crawler

import (
	"fmt"

	"github.com/PuerkitoBio/goquery"
)

//Current date and page postion related functions
func getLineOperationType(doc *goquery.Document) OperationType {
	operationTypeRaw := doc.Find(".schedule_active_list_active_tab").Text()

	operationType, ok := operationsIdentifiers[operationTypeRaw]
	if !ok {
		panic(fmt.Errorf("No operation mode found %v", operationTypeRaw))
	}
	return operationType
}
