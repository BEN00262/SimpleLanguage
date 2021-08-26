package main

func (parser *Parser) ParseIfStatement() interface{} {
	// expect the 'if' keyword
	parser.IsExpectedEatElsePanic(
		parser.CurrentToken(),
		KEYWORD, IF,
		"Expected the 'if' keyword",
	)

	condition := parser._parseExpression()

	parser.IsExpectedEatElsePanic(
		parser.CurrentToken(),
		CURLY_BRACES, "{",
		"Expected a '{'",
	)

	thenBody := parser._parseElseStatement()

	var elseBodies []interface{}

	// expect } --> closes the then block
	parser.IsExpectedEatElsePanic(
		parser.CurrentToken(),
		CURLY_BRACES, "}",
		"Expected a '}'",
	)

	if IsTypeAndValue(parser.CurrentToken(), KEYWORD, ELSE) {
		// eat the else keyword
		parser.eatToken()

		if IsTypeAndValue(parser.CurrentToken(), KEYWORD, IF) {
			elseBodies = append(elseBodies, parser.ParseIfStatement())
		} else {
			parser.IsExpectedEatElsePanic(
				parser.CurrentToken(),
				CURLY_BRACES, "{",
				"Expected a '{'",
			)

			elseBodies = append(elseBodies, BlockNode{
				Code: parser._parseElseStatement(),
			})

			parser.IsExpectedEatElsePanic(
				parser.CurrentToken(),
				CURLY_BRACES, "}",
				"Expected a '}'",
			)
		}

	}

	return IFNode{
		Condition: condition,
		ThenBody:  thenBody,
		ElseBody:  elseBodies,
	}
}

// can be also an expression
func (parser *Parser) _parseElseStatement() (statements []interface{}) {
	_currentToken := parser.CurrentToken()

	for !IsTypeAndValue(_currentToken, CURLY_BRACES, "}") {
		// loop
		statements = append(statements, parser._parse(_currentToken))
		_currentToken = parser.CurrentToken()
	}

	return
}
