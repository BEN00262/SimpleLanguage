package parser

import (
	"github.com/BEN00262/simpleLang/lexer"
)

func (parser *Parser) _parseBlockStatements() []interface{} {
	var _statements []interface{}
	_currentToken := parser.CurrentToken()

	for parser.CurrentPosition < parser.TokensLength && !IsTypeAndValue(_currentToken, lexer.CURLY_BRACES, "}") {
		_statements = append(_statements, parser._parse(_currentToken))
		_currentToken = parser.CurrentToken()
	}

	return _statements
}

func (parser *Parser) parseBlockScope() BlockNode {
	parser.IsExpectedEatElsePanic(
		parser.CurrentToken(),
		lexer.CURLY_BRACES, "{",
		"Expected '{'",
	)

	blockStatements := parser._parseBlockStatements()

	parser.IsExpectedEatElsePanic(
		parser.CurrentToken(),
		lexer.CURLY_BRACES, "}",
		"Expected '}'",
	)

	return BlockNode{
		Code: blockStatements,
	}
}
