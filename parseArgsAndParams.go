package main

type Param struct {
	Key      string
	Position int
}

func (parser *Parser) _parseFunctionParams() (params []Param, paramCount int) {
	_currentToken := parser.CurrentToken()
	parser.eatToken()

	_startPositon := 0

	for ; parser.CurrentPosition < parser.TokensLength && _currentToken.Type != HALF_CIRCLE_BRACKET; parser.eatToken() {

		if _currentToken.Type != VARIABLE {
			panic("Weird error")
		}

		params = append(params, Param{
			Key:      _currentToken.Value.(string),
			Position: _startPositon,
		})

		_startPositon += 1
		_currentToken = parser.CurrentToken()
	}

	parser.CurrentPosition -= 1

	paramCount = len(params)
	return
}

func (parser *Parser) _parseFunctionArgs() (_args []interface{}) {
	_currentToken := parser.CurrentToken()

	for parser.CurrentPosition < parser.TokensLength && !IsTypeAndValue(_currentToken, HALF_CIRCLE_BRACKET, ")") {
		_args = append(_args, parser._parseExpression())
		_currentToken = parser.CurrentToken()
	}

	return
}
