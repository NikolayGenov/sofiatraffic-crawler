package crawler

import "fmt"

type Operation int

type OperationID string

type OperationIDMap map[Operation]OperationID

type OperationIDRoutesMap map[OperationID]Routes

const (
	Operation_Normal Operation = iota
	Operation_Pre_Holiday
	Operation_Holiday
)

var (
	operationsIdentifiers = map[string]Operation{
		"делник":                         Operation_Normal,
		"предпразник":                    Operation_Pre_Holiday,
		"празник":                        Operation_Holiday,
		"предпразник / празник":          Operation_Holiday,
		"делник / предпразник / празник": Operation_Normal}

	operationStrings = [...]string{Operation_Normal: "Weekday",
		Operation_Pre_Holiday: "Pre-Holiday",
		Operation_Holiday:     "Holiday"}
)

func (o Operation) String() string {
	return operationStrings[o]
}

func convertToOperation(identifier string) (Operation, error) {
	t, ok := operationsIdentifiers[identifier]
	if !ok {
		return -1, fmt.Errorf("Unrecognized identifer for Operation type: %v", identifier)
	}
	return t, nil
}
