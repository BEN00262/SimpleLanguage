package parser

import (
	"fmt"

	. "github.com/BEN00262/simpleLang/lexer"
)

func (parser *Parser) _parseBody() []interface{} {
	var _body []interface{}
	_currentToken := parser.CurrentToken()

	for parser.CurrentPosition < parser.TokensLength && !IsTypeAndValue(_currentToken, CURLY_BRACES, "}") {
		_body = append(_body, parser._parse(_currentToken))
		_currentToken = parser.CurrentToken()
	}

	return _body
}

func (parser *Parser) _parseBlock() (block []interface{}) {
	// check for braces
	// then parse the blocks
	parser.IsExpectedEatElsePanic(
		parser.CurrentToken(),
		CURLY_BRACES, "{",
		fmt.Sprintf("Expected a '{' got a '%#v'", parser.CurrentToken().Value),
	)

	// start parsing the programs
	block = parser._parseBody() // error thrown from here will be passed down to the catch block and assigned to the name given ---> i think

	parser.IsExpectedEatElsePanic(
		parser.CurrentToken(),
		CURLY_BRACES, "}",
		fmt.Sprintf("Expected a '}' got a '%#v'", parser.CurrentToken().Value),
	)

	return
}

func (parser *Parser) parseTryCatchBlock() interface{} {
	parser.pushToParsingState(FUNCTION_STATE)
	defer parser.popFromParsingState()

	parser.IsExpectedEatElsePanic(
		parser.CurrentToken(),
		KEYWORD, TRY,
		fmt.Sprintf("Expected a 'try' statement got a '%#v'", parser.CurrentToken().Value),
	)

	// start parsing the programs
	tryBody := parser._parseBlock() // error thrown from here will be passed down to the catch block and assigned to the name given ---> i think

	parser.IsExpectedEatElsePanic(
		parser.CurrentToken(),
		KEYWORD, CATCH,
		fmt.Sprintf("Expected a 'catch' got a '%#v'", parser.CurrentToken().Value),
	)

	parser.IsExpectedEatElsePanic(
		parser.CurrentToken(),
		HALF_CIRCLE_BRACKET, "(",
		fmt.Sprintf("Expected a '(' got a '%#v'", parser.CurrentToken().Value),
	)

	_currentToken := parser.CurrentToken()

	// this is a parser error just panic --> we need to have states --> to work with this error stuffs
	if _currentToken.Type != VARIABLE {
		panic(fmt.Sprintf("Expected a type variable but got %#v", _currentToken.Value))
	}

	errorVariable := _currentToken.Value.(string)
	parser.eatToken()

	parser.IsExpectedEatElsePanic(
		parser.CurrentToken(),
		HALF_CIRCLE_BRACKET, ")",
		fmt.Sprintf("Expected a ')' got a '%#v'", parser.CurrentToken().Value),
	)

	// start parsing the programs
	catchBody := parser._parseBlock() // error thrown from here will be passed down to the catch block and assigned to the name given ---> i think

	// check if the next token is a finally block
	var finallyBody []interface{}

	if IsTypeAndValue(parser.CurrentToken(), KEYWORD, FINALLY) {
		parser.eatToken() // eat the finally keyword

		finallyBody = parser._parseBlock() // error thrown from here will be passed down to the catch block and assigned to the name given ---> i think
	}

	return TryCatchNode{
		Try: tryBody,
		Catch: CatchBlock{
			Exception: errorVariable,
			Body:      catchBody,
		},
		Finally: finallyBody,
	}
}
