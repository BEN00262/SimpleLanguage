package main

import "fmt"

type ContextValue = map[string]interface{}

type SymbolsTable struct {
	Contexts        []ContextValue
	CurrentPosition int
}

func initSymbolsTable() *SymbolsTable {
	return &SymbolsTable{
		CurrentPosition: -1,
	}
}

func (symbolsTable *SymbolsTable) pushContext() {
	symbolsTable.CurrentPosition += 1
	symbolsTable.Contexts = append(symbolsTable.Contexts, ContextValue{
		"_": nil,
	})
}

func (symbolsTable *SymbolsTable) popContext() {
	symbolsTable.Contexts = symbolsTable.Contexts[:len(symbolsTable.Contexts)-1]
	symbolsTable.CurrentPosition -= 1
}

func (symbolsTable *SymbolsTable) pushToContext(key string, value interface{}) {
	symbolsTable.Contexts[symbolsTable.CurrentPosition][key] = value
}

func (symbolsTable *SymbolsTable) getFromContext(key string) (interface{}, error) {
	lengthOfSymbolsTable := len(symbolsTable.Contexts) - 1

	for ; lengthOfSymbolsTable > -1; lengthOfSymbolsTable-- {
		currentContext := symbolsTable.Contexts[lengthOfSymbolsTable]

		for k, v := range currentContext {
			if k == key {
				return v, nil
			}
		}
	}

	return nil, fmt.Errorf("Varible '%s' does not exist", key)
}
