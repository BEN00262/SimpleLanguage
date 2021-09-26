package parser

import (
	"reflect"
	"testing"

	. "github.com/BEN00262/simpleLang/lexer"
)

func TestParsing(t *testing.T) {
	// we pass an array of tokens and then look
	tokens := []Token{
		{Type: KEYWORD, Value: "def"},
		{Type: VARIABLE, Value: "name"},
		{Type: ASSIGN, Value: "="},
		{Type: NUMBER, Value: "90"},
	}

	parser := InitParser(tokens)
	programNode := parser.Parse()

	t.Error(reflect.TypeOf(programNode))
}
