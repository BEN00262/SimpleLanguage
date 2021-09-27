package evaluator

import (
	"fmt"

	. "github.com/BEN00262/simpleLang/exceptions"
	. "github.com/BEN00262/simpleLang/parser"
)

func CheckFunctionArgs(function interface{}) {

}

func DaisyInvoke(eval *Evaluator, key string, args ...interface{}) (interface{}, ExceptionNode) {
	value, err := eval.symbolsTable.GetFromContext(key)

	if err != nil {
		panic(err)
	}

	if _value, ok := (*value).(SymbolTableValue); ok {
		if _value.Type == VALUE || _value.Type == ARRAY {
			return _value.Value, ExceptionNode{Type: NO_EXCEPTION}
		}

		eval.symbolsTable.PushContext()
		defer eval.symbolsTable.PopContext()

		// check the number of args passed
		// drop into a scope
		// inject the args into the scope
		// call the function

		switch _function_decl_ := _value.Value.(type) {
		case FunctionDecl:
			{
				if _function_decl_.ParamCount != len(args) {
					// panic
					panic(fmt.Sprintf("%s expected %d params but only %d args given", _function_decl_.Name, _function_decl_.ParamCount, len(args)))
				}

				// inject the args into the context
				for _, param := range _function_decl_.Params {
					valueType := VALUE

					_currentArg := args[param.Position]

					switch _currentArg.(type) {
					case AnonymousFunction:
						{
							valueType = FUNCTION
						}
					case ArrayNode:
						{
							valueType = ARRAY
						}
					}

					eval.symbolsTable.PushToContext(param.Key, SymbolTableValue{
						Type:  valueType,
						Value: _currentArg,
					})
				}

				fmt.Println("we are here")

				return eval.executeFunctionCode(_function_decl_.Code)
			}
		case AnonymousFunction:
			{
				return eval.executeFunctionCode(_function_decl_.Code)
			}
		}
	}

	return NilNode{}, ExceptionNode{Type: NO_EXCEPTION}
}
