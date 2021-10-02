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
	thenBody := parser.parseBlockScope()

	var elseBodies []interface{}

	if IsTypeAndValue(parser.CurrentToken(), KEYWORD, ELSE) {
		parser.eatToken() // else

		if IsTypeAndValue(parser.CurrentToken(), KEYWORD, IF) {
			elseBodies = append(elseBodies, parser.ParseIfStatement())
		} else {
			elseBodies = append(elseBodies, parser.parseBlockScope())
		}
	}

	return IFNode{
		Condition: condition,
		ThenBody:  thenBody,
		ElseBody:  elseBodies,
	}
}
