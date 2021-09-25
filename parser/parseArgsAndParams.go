package parser

import (
	. "github.com/BEN00262/simpleLang/lexer"
)

type Param struct {
	Key      string
	Position int
}

func (parser *Parser) _parseFunctionParams() (params []Param, paramCount int) {
	_currentToken := parser.CurrentToken()
	_startPositon := 0

	for parser.CurrentPosition < parser.TokensLength && !IsTypeAndValue(_currentToken, HALF_CIRCLE_BRACKET, ")") {
		if _currentToken.Type != VARIABLE {
			panic("Weird error")
		}

		params = append(params, Param{
			Key:      _currentToken.Value.(string),
			Position: _startPositon,
		})

		_startPositon += 1
		parser.eatToken()

		if IsTypeAndValue(parser.CurrentToken(), COMMA, ",") {
			parser.eatToken()
		}

		_currentToken = parser.CurrentToken()
	}

	paramCount = len(params)
	return
}

func (parser *Parser) _parseFunctionArgs() (_args []interface{}) {
	_currentToken := parser.CurrentToken()

	for parser.CurrentPosition < parser.TokensLength && !IsTypeAndValue(_currentToken, HALF_CIRCLE_BRACKET, ")") {
		_args = append(_args, parser._parseExpression())
		_currentToken = parser.CurrentToken()

		if _currentToken.Type == COMMA {
			// eat the damn shit
			parser.eatToken()
		}
	}

	return
}
