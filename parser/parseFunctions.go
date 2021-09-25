package parser

import (
	"fmt"

	. "github.com/BEN00262/simpleLang/lexer"
)

func IsTypeAndValue(token Token, expectedType TokenType, value string) bool {
	return token.Type == expectedType && token.Value.(string) == value
}

// this eats the next token
func (parser *Parser) IsExpectedEatElsePanic(token Token, expectedType TokenType, value string, panicMessage string) {
	if !IsTypeAndValue(token, expectedType, value) {
		panic(panicMessage)
	}
	parser.eatToken() // advance the counter if true
}

func (parser *Parser) parseCommonFunctionCode() ([]interface{}, []Param, int) {
	// eat the (
	parser.IsExpectedEatElsePanic(
		parser.CurrentToken(),
		HALF_CIRCLE_BRACKET, "(",
		fmt.Sprintf("Expected '(' but got '%s'", parser.CurrentToken().Value),
	)

	// parse the params here
	params, paramCount := parser._parseFunctionParams()

	// eat the )
	parser.IsExpectedEatElsePanic(
		parser.CurrentToken(),
		HALF_CIRCLE_BRACKET, ")",
		fmt.Sprintf("Expected ')' but got '%s'", parser.CurrentToken().Value),
	)

	// eat the {
	parser.IsExpectedEatElsePanic(
		parser.CurrentToken(),
		CURLY_BRACES, "{",
		fmt.Sprintf("Expected '{' but got '%s'", parser.CurrentToken().Value),
	)

	_currentToken := parser.CurrentToken()
	var _code []interface{}

	for !IsTypeAndValue(_currentToken, CURLY_BRACES, "}") {
		_code = append(_code, parser._parse(_currentToken))
		_currentToken = parser.CurrentToken()
	}

	// eat the {
	parser.IsExpectedEatElsePanic(
		parser.CurrentToken(),
		CURLY_BRACES, "}",
		fmt.Sprintf("Expected '}' but got '%s'", parser.CurrentToken().Value),
	)

	return _code, params, paramCount
}

func (parser *Parser) ParseFunction() interface{} {
	// enter the function_state inorder to validate keywords like return
	parser.pushToParsingState(FUNCTION_STATE)
	defer parser.popFromParsingState()

	_func_ := parser.CurrentToken()

	parser.IsExpectedEatElsePanic(
		_func_,
		KEYWORD, FUNC,
		"Not a valid function declaration expected 'fun' keyword",
	)

	_current_token_ := parser.CurrentToken()

	if IsTypeAndValue(_current_token_, HALF_CIRCLE_BRACKET, "(") {
		_code, params, paramCount := parser.parseCommonFunctionCode()

		return AnonymousFunction{
			ParamCount: paramCount,
			Params:     params,
			Code:       _code,
		}
	}

	_name := _current_token_.Value.(string) // function name
	parser.eatToken()                       // proceed to the (

	_code, params, paramCount := parser.parseCommonFunctionCode()

	return FunctionDecl{
		Name:       _name,
		ParamCount: paramCount,
		Params:     params,
		Code:       _code,
	}
}
