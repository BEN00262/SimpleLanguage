package parser

import (
	. "github.com/BEN00262/simpleLang/lexer"
)

func (parser *Parser) ParseAssignment() (lvalue interface{}) {
	if IsTypeAndValue(parser.CurrentToken(), KEYWORD, FUNC) {
		lvalue = parser.ParseFunction()
	} else {
		lvalue = parser._parseExpression()
	}
	return
}
