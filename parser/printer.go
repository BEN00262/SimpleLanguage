package parser

import (
	"fmt"
	"strings"
)

func Print(value interface{}) string {
	switch _value := value.(type) {
	case StringNode:
		{
			return fmt.Sprintf("%s", _value.Value)
		}
	case NumberNode:
		{
			if _value.Type == INTEGER {
				return fmt.Sprintf("%v", _value.Value.Text(10))
			}

			return fmt.Sprintf("%v", _value.FValue.Text('f', 64))
		}
	case BoolNode:
		{
			if _value.Value == 1 {
				return fmt.Sprint("True")
			} else {
				return fmt.Sprint("False")
			}
		}
	case ArrayNode:
		{
			var _arguments_ []string

			for _, _element := range _value.Elements {
				_arguments_ = append(_arguments_, Print(_element))
			}

			return "[ " + strings.Join(_arguments_, ",") + " ]"
		}
	}
	return ""
}
