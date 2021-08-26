package main

import "fmt"

func (parser *Parser) _parseForLoop() []interface{} {
	var _whileBody []interface{}
	_currentToken := parser.CurrentToken() // we start with the actual thing

	for ; parser.CurrentPosition < parser.TokensLength && !IsTypeAndValue(_currentToken, CURLY_BRACES, "}"); parser.CurrentPosition++ {
		_parsed_ := parser._parse(_currentToken)

		fmt.Println("After parsing")
		fmt.Println(_parsed_)
		fmt.Println(parser.CurrentToken())

		_whileBody = append(_whileBody, _parsed_)
		_currentToken = parser.CurrentToken()

		fmt.Println("Current tokens being worked on")
		fmt.Println(_currentToken)
		fmt.Println(parser.peekAhead())
	}

	parser.CurrentPosition -= 1 // one for the while loop and the other increment for the shit in parseExpression

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
	fmt.Println("Tuko hapa")
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

	_lookAheadToken := parser.peekAhead()

	if IsTypeAndValue(_lookAheadToken, ASSIGN, "=") {
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
