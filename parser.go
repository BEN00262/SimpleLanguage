package main

import (
	"strings"
)

type Parser struct {
	Tokens          []Token
	CurrentPosition int
	TokensLength    int
}

func initParser(tokens []Token) *Parser {
	return &Parser{
		Tokens:          tokens,
		CurrentPosition: 0,
		TokensLength:    len(tokens),
	}
}

func (parser *Parser) CurrentToken() Token {
	if parser.CurrentPosition < parser.TokensLength {
		return parser.Tokens[parser.CurrentPosition]
	}

	return Token{}
}

func (parser *Parser) peekAhead() Token {
	futurePosition := parser.CurrentPosition + 1

	if futurePosition < parser.TokensLength {
		return parser.Tokens[futurePosition]
	}

	return Token{}
}

func (parser *Parser) eatToken() {
	parser.CurrentPosition++
}

func (parser *Parser) _parseFactorTail() interface{} {
	// parser.eatToken()
	_currentToken := parser.CurrentToken()

	if _currentToken.Type == OPERATOR && strings.Contains("*/", _currentToken.Value.(string)) {
		// check if we are here

		parser.eatToken()

		factor := parser._parseFactor()
		// parser.eatToken()
		factorTail := parser._parseFactorTail()

		if factorTail != nil {
			_factorTail := factorTail.(BinaryNode)
			_factorTail.Lhs = factor
			factor = _factorTail
		}

		return BinaryNode{
			Rhs:      factor,
			Operator: _currentToken.Value.(string),
		}
	}

	return nil
}

func (parser *Parser) _parseFactor() interface{} {
	_currentToken := parser.CurrentToken()

	// check the currentToken for the value
	if _currentToken.Type == HALF_CIRCLE_BRACKET && _currentToken.Value == '(' {

		parser.eatToken()

		_expression := parser._parseExpression()

		_currentToken = parser.CurrentToken()

		if _currentToken.Type != HALF_CIRCLE_BRACKET && _currentToken.Value != ')' {
			// throw an error here
		}

		parser.eatToken()
		return _expression
	} else if _currentToken.Type == NUMBER {
		return NumberNode{
			Value: _currentToken.Value.(int),
		}

	} else if _currentToken.Type == VARIABLE {
		_peakToken := parser.peekAhead()

		if _peakToken.Type == HALF_CIRCLE_BRACKET && _peakToken.Value == '(' {
			// eat the variable
			parser.eatToken()
			parser.eatToken()

			_currenToken := parser.CurrentToken()
			var _array_of_args []interface{}

			for _currenToken.Type != HALF_CIRCLE_BRACKET {
				if _currenToken.Type == COMMA {
					parser.eatToken()
					_currenToken = parser.CurrentToken()
					continue
				}

				_result := parser._parseExpression()
				_array_of_args = append(_array_of_args, _result)

				parser.eatToken()
				_currenToken = parser.CurrentToken()
			}

			// we loop till we get to the ) bracket

			// evaluate the expression in the args section then return the function call Node
			if _currentToken.Type != HALF_CIRCLE_BRACKET && _currentToken.Value != ')' {
				// throw an error here
			}

			parser.eatToken()

			// get the args here ( start the operation of using them )
			return FunctionCall{
				Name:     _currentToken.Value.(string),
				ArgCount: 0,
			}

		}

		return VariableNode{
			Value: _currentToken.Value.(string),
		}
	}

	return nil
}

func (parser *Parser) _parseTermTail() interface{} {
	_currentToken := parser.CurrentToken()

	if _currentToken.Type == OPERATOR && strings.Contains("+-", _currentToken.Value.(string)) {
		parser.eatToken()

		term := parser._parseTerm()
		termTail := parser._parseTermTail()

		if termTail != nil {
			_termTail := termTail.(BinaryNode)
			_termTail.Lhs = term
			term = _termTail
		}

		return BinaryNode{
			Rhs:      term,
			Operator: _currentToken.Value.(string),
		}
	}

	return nil
}

func (parser *Parser) _parseTerm() interface{} {
	factor := parser._parseFactor()

	parser.eatToken()
	factorTail := parser._parseFactorTail()
	// parser.eatToken()

	// check to ensure that the factor Tail is not nil if so just return the factor

	if factorTail != nil {
		_factorTail := factorTail.(BinaryNode)
		_factorTail.Lhs = factor
		factor = _factorTail
	}

	return factor
}

func (parser *Parser) _parseExpression() interface{} {
	term := parser._parseTerm()

	termTail := parser._parseTermTail()

	if termTail != nil {
		_termTail := termTail.(BinaryNode)
		_termTail.Lhs = term
		term = _termTail
	}

	return ExpressionNode{
		expression: term,
	}
}

func (parser *Parser) _parse(token Token) interface{} {
	switch token.Type {
	case KEYWORD:
		{
			if token.Value == FUNC {
				parser.eatToken()

				_funcName := parser.CurrentToken()

				if _funcName.Type != VARIABLE {
					// raise an error here
					// panic here
				}

				_name := _funcName.Value.(string)

				parser.eatToken()

				if parser.peekAhead().Type != HALF_CIRCLE_BRACKET {
					// raise an error here
				}

				parser.eatToken()

				if parser.peekAhead().Type != HALF_CIRCLE_BRACKET {
					// raise and error here
				}

				parser.eatToken()

				if parser.peekAhead().Type != CURLY_BRACES {
					// raise another error here
				}

				parser.eatToken()

				_currentToken := parser.CurrentToken()
				var _code []interface{}

				for _currentToken.Type != CURLY_BRACES {
					_code = append(_code, parser._parse(_currentToken))
					_currentToken = parser.CurrentToken()
				}

				if parser.CurrentToken().Type != CURLY_BRACES {
					// raise another error here
				}

				// fmt.Println("Last token from the parser")
				// fmt.Println(parser.CurrentToken())

				return FunctionDecl{
					Name:       _name,
					ParamCount: 0,
					Code:       _code,
				}
			}
		}
	case VARIABLE:
		{
			nextToken := parser.peekAhead()

			if nextToken.Type == ASSIGN {
				// eat the var too
				parser.eatToken()
				parser.eatToken()

				lvalue := parser._parseExpression()

				return Assignment{
					Lvalue: token.Value.(string),
					Rvalue: lvalue,
				}

			}

			return parser._parseExpression()
		}

	case COMMENT:
		{
			// place code in this comments but how are we to parse them
			// made allow interpolation in the comments who knows
			// do the literate programming here
			// extract the code from the comments then execute it
			return CommentNode{
				comment: token.Value.(string),
			}
		}

	default:
		{
			return parser._parseExpression()
		}
	}

	return nil
}

func (parser *Parser) Parse() *ProgramNode {
	program := &ProgramNode{}

	for parser.CurrentPosition < parser.TokensLength {
		program.Nodes = append(program.Nodes, parser._parse(parser.CurrentToken()))
		parser.CurrentPosition += 1
	}

	return program
}
