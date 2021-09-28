package parser

import (
	. "github.com/BEN00262/simpleLang/lexer"
)

// parse function call
func (parser *Parser) parseFunctionCall(objectAccessor ObjectAccessor) interface{} {
	parser.IsExpectedEatElsePanic(
		parser.CurrentToken(),
		HALF_CIRCLE_BRACKET, "(",
		"Expectd '('",
	)

	_array_of_args := parser._parseFunctionArgs()

	parser.IsExpectedEatElsePanic(
		parser.CurrentToken(),
		HALF_CIRCLE_BRACKET, ")",
		"Expected a ')'",
	)

	return FunctionCall{
		Name:     objectAccessor,
		ArgCount: len(_array_of_args),
		Args:     _array_of_args,
	}
}

func (parser *Parser) parseArrayIndexing(objectAccessor ObjectAccessor) interface{} {
	parser.IsExpectedEatElsePanic(
		parser.CurrentToken(),
		SQUARE_BRACKET, "[",
		"Expectd '['",
	)

	// parse the expression here
	_accessor_type_ := NORMAL
	var _end_index_expression_ interface{}
	_array_index_expression_ := parser._parseExpression()

	// check if the current token is a :

	// if so eat it and do shit
	if IsTypeAndValue(parser.CurrentToken(), COLON, ":") {
		parser.eatToken()

		// we should check for errors later

		_end_index_expression_ = parser._parseExpression()
		_accessor_type_ = RANGE
	}

	parser.IsExpectedEatElsePanic(
		parser.CurrentToken(),
		SQUARE_BRACKET, "]",
		"Expected a closing ']'",
	)

	return ArrayAccessorNode{
		Array:    objectAccessor,
		Index:    _array_index_expression_,
		Type:     _accessor_type_,
		EndIndex: _end_index_expression_,
	}
}

func (parser *Parser) parseVariableExpression() interface{} {
	_currentToken := parser.CurrentToken()
	parser.eatToken()

	if _currentToken.Type != VARIABLE {
		panic("We expected a variable")
	}

	if IsTypeAndValue(parser.CurrentToken(), DOT, ".") {
		parser.IsExpectedEatElsePanic(
			parser.CurrentToken(),
			DOT, ".",
			"Expected a '.'",
		)

		// get the next token
		_child := parser.CurrentToken()

		if _child.Type != VARIABLE {
			// we should panic here
			panic("Expected a variable but got something else")
		}

		parser.eatToken()

		object := ObjectAccessor{
			Parent: _currentToken.Value.(string),
			Child:  _child.Value.(string),
		}

		if IsTypeAndValue(parser.CurrentToken(), SQUARE_BRACKET, "[") {
			return parser.parseArrayIndexing(object)
		}

		if IsTypeAndValue(parser.CurrentToken(), HALF_CIRCLE_BRACKET, "(") {
			return parser.parseFunctionCall(object)
		}

		return object
	}

	// this is the other stuff right here
	object := ObjectAccessor{Parent: _currentToken.Value.(string)}

	if IsTypeAndValue(parser.CurrentToken(), SQUARE_BRACKET, "[") {
		return parser.parseArrayIndexing(object)
	}

	if IsTypeAndValue(parser.CurrentToken(), HALF_CIRCLE_BRACKET, "(") {
		return parser.parseFunctionCall(object)
	}

	return VariableNode{
		Value: _currentToken.Value.(string),
	}
}
