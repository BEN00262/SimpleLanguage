package main

import "fmt"

type ExternalDependencies = map[string]SymbolTableValue

var (
	GLOBALS = ExternalDependencies{
		"print": SymbolTableValue{
			Type: EXTERNALFUNC,
			Value: ExternalFunctionNode{
				Name:       "print",
				ParamCount: 1,
				Function: func(values ...interface{}) interface{} {
					for _, value := range values {
						switch _value := value.(type) {
						case StringNode:
							{
								fmt.Printf("%s", _value.Value)
							}
						case NumberNode:
							{
								fmt.Printf("%d", _value.Value)
							}
						}
					}

					fmt.Println()
					return NilNode{}
				},
			},
		},
	}
)

func LoadGlobalsToContext(eval *Evaluator) {
	for Key, Value := range GLOBALS {
		eval.InjectIntoGlobalScope(Key, Value)
	}
}
