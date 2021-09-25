package evaluator

import (
	"reflect"

	. "github.com/BEN00262/simpleLang/parser"
)

func ToDaisy(value interface{}) interface{} {
	_valueKind := reflect.ValueOf(value)

	switch _valueKind.Kind() {
	case reflect.Int:
		{
			return NumberNode{
				Value: int(_valueKind.Int()),
			}
		}
	case reflect.String:
		{
			return StringNode{
				Value: _valueKind.String(),
			}
		}
	}
	return NilNode{}
}

func FromDaisy(value interface{}) interface{} {
	// get the values and convert them to golang types
	switch _val := value.(type) {
	case NumberNode:
		{
			return _val.Value
		}
	case StringNode:
		{
			return _val.Value
		}
	}

	return nil
}
