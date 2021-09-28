package parser

import (
	. "github.com/BEN00262/simpleLang/lexer"
)

func (parser *Parser) ParseIfStatement() interface{} {
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

	thenBody := parser._parseForLoop()

	var elseBodies []interface{}

	parser.IsExpectedEatElsePanic(
		parser.CurrentToken(),
		CURLY_BRACES, "}",
		"Expected a '}'",
	)

	if IsTypeAndValue(parser.CurrentToken(), KEYWORD, ELSE) {
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
				Code: parser._parseForLoop(),
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
