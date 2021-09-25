package parser

import (
	"fmt"

	. "github.com/BEN00262/simpleLang/lexer"
)

func (parser *Parser) _parseForLoop() []interface{} {
	var _whileBody []interface{}
	_currentToken := parser.CurrentToken()

	// get the current Token --> proceed --> the problem starts with the expressions

	for parser.CurrentPosition < parser.TokensLength && !IsTypeAndValue(_currentToken, CURLY_BRACES, "}") {
		_whileBody = append(_whileBody, parser._parse(_currentToken))
		_currentToken = parser.CurrentToken()
	}

	// parser.CurrentPosition -= 1 // one for the while loop and the other increment for the shit in parseExpression
	return _whileBody
}

func (parser *Parser) parseWhileLoop() interface{} {
	// expect {
	parser.IsExpectedEatElsePanic(
		parser.CurrentToken(),
		CURLY_BRACES, "{",
		fmt.Sprintf("Expected '{' but got a '%#v'", parser.CurrentToken().Value),
	)

	_for_body_ := parser._parseForLoop()

	// expect }
	parser.IsExpectedEatElsePanic(
		parser.CurrentToken(),
		CURLY_BRACES, "}",
		fmt.Sprintf("Expected '}' but got a '%#v'", parser.CurrentToken().Value),
	)

	return ForNode{
		Type:    WHILE_FOREVER,
		ForBody: _for_body_,
	}
}

func (parser *Parser) ParseForLoop() interface{} {
	// eat the for after asserting it presence
	parser.IsExpectedEatElsePanic(
		parser.CurrentToken(),
		KEYWORD, FOR,
		"Expected the 'for' keyword",
	)

	_currentToken := parser.CurrentToken()

	// the while portion should gurantee shit here
	if IsTypeAndValue(_currentToken, CURLY_BRACES, "{") {
		return parser.parseWhileLoop()
	}

	// expect (
	parser.IsExpectedEatElsePanic(
		_currentToken,
		HALF_CIRCLE_BRACKET, "(",
		"Expected a '('",
	)

	// check if there is an eassignment thing if not assume its a conditional while loop

	if IsTypeAndValue(parser.CurrentToken(), KEYWORD, "def") {
		initialization := parser._parse(parser.CurrentToken())

		// expect ;
		parser.IsExpectedEatElsePanic(
			parser.CurrentToken(),
			SEMI_COLON, ";",
			"Expected a ';'",
		)

		conditional := parser._parseExpression()

		// expect ;
		parser.IsExpectedEatElsePanic(
			parser.CurrentToken(),
			SEMI_COLON, ";",
			"Expected a ';'",
		)

		increment := parser._parseExpression()

		// expect )
		parser.IsExpectedEatElsePanic(
			parser.CurrentToken(),
			HALF_CIRCLE_BRACKET, ")",
			"Expected a ')'",
		)

		// expect {
		parser.IsExpectedEatElsePanic(
			parser.CurrentToken(),
			CURLY_BRACES, "{",
			"Expected a '{'",
		)

		_for_body_ := parser._parseForLoop()

		// expect }
		parser.IsExpectedEatElsePanic(
			parser.CurrentToken(),
			CURLY_BRACES, "}",
			"Expected a '}'",
		)

		return ForNode{
			Type:           FOR_NODE,
			Initialization: initialization,
			Condition:      conditional,
			Increment:      increment,
			ForBody:        _for_body_,
		}

	}

	// conditional loop

	// parse the conditional
	conditional := parser._parseExpression()

	// expect )
	parser.IsExpectedEatElsePanic(
		parser.CurrentToken(),
		HALF_CIRCLE_BRACKET, ")",
		"Expected a ')'",
	)

	// expect {
	parser.IsExpectedEatElsePanic(
		parser.CurrentToken(),
		CURLY_BRACES, "{",
		"Expected a '{'",
	)

	_forBody := parser._parseForLoop()

	// expect }
	parser.IsExpectedEatElsePanic(
		parser.CurrentToken(),
		CURLY_BRACES, "}",
		"Expected a '}'",
	)

	return ForNode{
		Type:      WHILE_CONDITIONAL,
		Condition: conditional,
		ForBody:   _forBody,
	}
}
