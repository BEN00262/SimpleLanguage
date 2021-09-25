package symbolstable

import "fmt"

type ContextValue = map[string]interface{}

type SymbolsTable struct {
	Contexts        []ContextValue
	CurrentPosition int
}

func InitSymbolsTable() *SymbolsTable {
	return &SymbolsTable{
		CurrentPosition: -1,
	}
}

func (symbolsTable *SymbolsTable) PushContext() {
	symbolsTable.CurrentPosition += 1
	symbolsTable.Contexts = append(symbolsTable.Contexts, ContextValue{
		"_": nil,
	})
}

func (symbolsTable *SymbolsTable) PopContext() {
	symbolsTable.Contexts = symbolsTable.Contexts[:len(symbolsTable.Contexts)-1]
	symbolsTable.CurrentPosition -= 1
}

func (symbolsTable *SymbolsTable) PushToContext(key string, value interface{}) {
	symbolsTable.Contexts[symbolsTable.CurrentPosition][key] = value
}

// we need to find a position in the context and push the value in there
func (symbolsTable *SymbolsTable) PushToParentContext(key string, value interface{}) error {
	lengthOfSymbolsTable := len(symbolsTable.Contexts) - 1

	for ; lengthOfSymbolsTable > -1; lengthOfSymbolsTable-- {
		currentContext := symbolsTable.Contexts[lengthOfSymbolsTable]

		for k := range currentContext {
			if k == key {
				currentContext[k] = value
				return nil
			}
		}
	}

	return fmt.Errorf("Varible '%s' does not exist", key)
}

func (symbolsTable *SymbolsTable) GetFromContext(key string) (*interface{}, error) {
	lengthOfSymbolsTable := len(symbolsTable.Contexts) - 1

	// return pointers to things so that we can change them up
	for ; lengthOfSymbolsTable > -1; lengthOfSymbolsTable-- {
		currentContext := symbolsTable.Contexts[lengthOfSymbolsTable]

		for k, v := range currentContext {
			if k == key {
				return &v, nil
			}
		}
	}

	return nil, fmt.Errorf("Varible '%s' does not exist", key)
}