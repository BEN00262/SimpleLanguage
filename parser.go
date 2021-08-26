package main

import (
	"strings"
)

type ParsingState = int

const (
	FUNCTION_STATE ParsingState = iota
	LOOP_STATE
	INVALID_STATE // used when we are not in a function or a loop
)

// how the fuck will we do this
type Parser struct {
	Tokens          []Token
	State           []ParsingState // push into this state
	CurrentPosition int
	TokensLength    int
}

func initParser(tokens []Token) *Parser {
	return &Parser{
		Tokens:          tokens,
		CurrentPosition: 0,
		State:           []ParsingState{INVALID_STATE},
		TokensLength:    len(tokens),
	}
}

func (parser *Parser) pushToParsingState(state ParsingState) {
	parser.State = append(parser.State, state)
}

func (parser *Parser) popFromParsingState() {
	parser.State = parser.State[0 : len(parser.State)-1]
}

// we have to think we cannot have a break in a function alone
// so in break the top state must be loo_state
// for return the parent state must be the function state
func (parser *Parser) EnsureCurrentParsingStateIs(state ParsingState, dig bool) bool {
	if dig {
		// start from the end going to the front trying to check if we are in a function or what
		// we need to make sure that the CurrentCounter is greater than 1 else just return the one
		for CurrentStateCounter := len(parser.State) - 1; CurrentStateCounter > -1; CurrentStateCounter-- {
			if parser.State[CurrentStateCounter] == state {
				return true
			}
		}

		return false
	}

	return parser.State[len(parser.State)-1] == state
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
	_currentToken := parser.CurrentToken()

	if _currentToken.Type == OPERATOR && strings.Contains("*/", _currentToken.Value.(string)) {
		parser.eatToken() // eat the operator * or /

		factor := parser._parseFactor()
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

	if IsTypeAndValue(_currentToken, HALF_CIRCLE_BRACKET, "(") {
		parser.eatToken() // eat the (

		if IsTypeAndValue(parser.CurrentToken(), KEYWORD, FUNC) {
			_iife_func_ := parser.ParseFunction()

			parser.IsExpectedEatElsePanic(
				parser.CurrentToken(),
				HALF_CIRCLE_BRACKET, ")",
				"Expected a ')'",
			)

			// we check if there is anything like a ( if there is parse it as an iife else parse it
			// as an anonymous function
			if !IsTypeAndValue(parser.CurrentToken(), HALF_CIRCLE_BRACKET, "(") {
				// return the func declaration as an anonymous thing
				return _iife_func_
			}

			// just incase its an iife
			_iife_function_ := _iife_func_.(AnonymousFunction)

			parser.IsExpectedEatElsePanic(
				parser.CurrentToken(),
				HALF_CIRCLE_BRACKET, "(",
				"Expected a '('",
			)

			_args := parser._parseFunctionArgs()

			parser.IsExpectedEatElsePanic(
				parser.CurrentToken(),
				HALF_CIRCLE_BRACKET, ")",
				"Expected a ')'",
			)

			return IIFENode{
				Function: _iife_function_,
				Args:     _args,
				ArgCount: len(_args),
			}
		}

		_expression := parser._parseExpression()

		parser.IsExpectedEatElsePanic(
			parser.CurrentToken(),
			HALF_CIRCLE_BRACKET, ")",
			"Expected a ')'",
		)

		return _expression
	} else if _currentToken.Type == NUMBER {
		parser.eatToken()

		return NumberNode{
			Value: _currentToken.Value.(int),
		}
	} else if _currentToken.Type == STRING {
		parser.eatToken()

		return StringNode{
			Value: _currentToken.Value.(string),
		}
	} else if _currentToken.Type == KEYWORD {
		parser.eatToken()

		if strings.Compare(_currentToken.Value.(string), TRUE) == 0 {
			return BoolNode{
				Value: 1,
			}
		}

		if strings.Compare(_currentToken.Value.(string), FALSE) == 0 {
			return BoolNode{
				Value: 0,
			}
		}

	} else if _currentToken.Type == VARIABLE {
		if IsTypeAndValue(parser.peekAhead(), HALF_CIRCLE_BRACKET, "(") {
			parser.eatToken() // function name
			parser.eatToken() // the first (

			_array_of_args := parser._parseFunctionArgs()

			parser.IsExpectedEatElsePanic(
				parser.CurrentToken(),
				HALF_CIRCLE_BRACKET, ")",
				"Expected a ')'",
			)

			return FunctionCall{
				Name:     _currentToken.Value.(string),
				ArgCount: len(_array_of_args),
				Args:     _array_of_args,
			}

		}

		// eat its own shit
		parser.eatToken()
		return VariableNode{
			Value: _currentToken.Value.(string),
		}
	}

	return nil
}

func (parser *Parser) _parseTermTail() interface{} {
	_currentToken := parser.CurrentToken()

	// also check for conditionals here
	if _currentToken.Type == OPERATOR && strings.Contains("+-", _currentToken.Value.(string)) {
		parser.eatToken() // eat the operator e.g + or -

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
	factorTail := parser._parseFactorTail()

	if factorTail != nil {
		_factorTail := factorTail.(BinaryNode)
		_factorTail.Lhs = factor
		factor = _factorTail
	}

	return factor
}

func (parser *Parser) _parseExpression() interface{} {
	_lhs := parser._parseComparison()

	_currentToken := parser.CurrentToken()

	// this fuckers should eat their shit
	if _currentToken.Type == CONDITION && Contains([]string{"!=", "=="}, _currentToken.Value.(string)) {
		parser.eatToken()

		_expression := parser._parseExpression()

		return ConditionNode{
			Operator: _currentToken.Value.(string),
			Rhs:      _expression,
			Lhs:      _lhs,
		}
	}

	return _lhs
}

func (parser *Parser) _parseInnerComparison() interface{} {
	term := parser._parseTerm()
	termTail := parser._parseTermTail()

	// this iss weird

	if termTail != nil {
		_termTail := termTail.(BinaryNode)
		_termTail.Lhs = term
		term = _termTail
	}

	return ExpressionNode{
		expression: term,
	}
}

func (parser *Parser) _parseComparison() interface{} {
	_lhs := parser._parseInnerComparison()
	_currentToken := parser.CurrentToken()

	if _currentToken.Type == CONDITION && Contains([]string{">", ">=", "<", "<="}, _currentToken.Value.(string)) {
		parser.eatToken()
		_rhs := parser._parseInnerComparison()

		return ConditionNode{
			Lhs:      _lhs,
			Rhs:      _rhs,
			Operator: _currentToken.Value.(string),
		}
	}

	return _lhs
}

// use a stack to monitor the current state we are --> are we parsing a function or a loop
// push the current parsing state to a stack then check on special tokens like return and break
func (parser *Parser) _parse(token Token) interface{} {
	switch token.Type {
	case KEYWORD:
		{
			switch token.Value {
			case FUNC:
				{
					return parser.ParseFunction()
				}
			case FOR:
				{
					// parser.eatToken()

					parser.pushToParsingState(LOOP_STATE)
					_for_node_ := parser.ParseForLoop()
					parser.popFromParsingState()

					return _for_node_
				}
			case IF:
				{
					// this parses an if statment
					return parser.ParseIfStatement()
				}
			case BREAK:
				{
					// check the current state we are in
					if !parser.EnsureCurrentParsingStateIs(LOOP_STATE, false) {
						// panic here
						panic("Break cannot be used outside a loop")
					}

					return BreakNode{}
				}
			case RETURN:
				{
					if !parser.EnsureCurrentParsingStateIs(FUNCTION_STATE, true) {
						// panic here
						panic("Return cannot be used outside a function definition")
					}

					parser.eatToken()
					_expression := parser._parseExpression()

					return ReturnNode{
						Expression: _expression,
					}
				}
			default:
				{
					return parser._parseExpression()
				}
			}
		}
	case VARIABLE:
		{
			if IsTypeAndValue(parser.peekAhead(), ASSIGN, "=") {
				parser.eatToken() // eat the variable
				parser.eatToken() // eat the assignment operator

				var lvalue interface{}

				// this is not right at all
				// the functions should be an expression
				if IsTypeAndValue(parser.CurrentToken(), KEYWORD, FUNC) {
					// parsing anonymous functions
					lvalue = parser.ParseFunction()
				} else {
					// expressions
					lvalue = parser._parseExpression()
				}

				return Assignment{
					Lvalue: token.Value.(string),
					Rvalue: lvalue,
				}

			}

			// why the fuck does this not execute
			return parser._parseExpression()
		}

	case COMMENT:
		{
			// increment the fuck out of this
			parser.eatToken()

			return CommentNode{
				comment: token.Value.(string),
			}
		}
	default:
		{
			return parser._parseExpression()
		}
	}
}

func (parser *Parser) Parse() *ProgramNode {
	program := &ProgramNode{}

	// start parsing the tokens here make everything to do its own complete shit
	for parser.CurrentPosition < parser.TokensLength {
		program.Nodes = append(program.Nodes, parser._parse(parser.CurrentToken()))
		// parser.CurrentPosition += 1
	}

	return program
}
