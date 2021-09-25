package parser

import (
	"strings"

	. "github.com/BEN00262/simpleLang/lexer"
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

func InitParser(tokens []Token) *Parser {
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

		// parse an array here and then return its type and stuff
		// we have a shit stuff here boys
	} else if IsTypeAndValue(_currentToken, SQUARE_BRACKET, "[") {
		// start parsing the array here
		// consume the first square bracke
		parser.eatToken() // eat the [

		var _elements_ []interface{}

		for parser.CurrentPosition < parser.TokensLength && !IsTypeAndValue(parser.CurrentToken(), SQUARE_BRACKET, "]") {
			_elements_ = append(_elements_, parser._parseExpression())

			if parser.CurrentToken().Type == COMMA {
				parser.eatToken()
			}
		}

		// final we expect a closing bracket
		// eat ]
		parser.IsExpectedEatElsePanic(
			parser.CurrentToken(),
			SQUARE_BRACKET, "]",
			"Expected ']'",
		)

		return ArrayNode{
			Elements: _elements_,
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

		if _currentToken.Value == NIL {
			return NilNode{}
		}

	} else if _currentToken.Type == VARIABLE {
		if IsTypeAndValue(parser.peekAhead(), HALF_CIRCLE_BRACKET, "(") {
			parser.eatToken()
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

		} else if IsTypeAndValue(parser.peekAhead(), SQUARE_BRACKET, "[") {
			// also check if the next token is a square bracket if so this is an array access thing
			parser.eatToken()
			parser.eatToken()

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
				Array:    _currentToken.Value.(string),
				Index:    _array_index_expression_,
				Type:     _accessor_type_,
				EndIndex: _end_index_expression_,
			}
		}

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
		Expression: term,
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
					parser.pushToParsingState(LOOP_STATE)
					_for_node_ := parser.ParseForLoop()
					parser.popFromParsingState()

					return _for_node_
				}
			case IF:
				{
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
			case DEF:
				{
					parser.eatToken() // eat the def keyword

					lvalue := parser.CurrentToken().Value.(string)

					parser.eatToken() // eat the variable name

					parser.IsExpectedEatElsePanic(
						parser.CurrentToken(),
						ASSIGN, "=",
						"Expected '='",
					)

					rvalue := parser.ParseAssignment()

					return Assignment{
						Type:   ASSIGNMENT,
						Lvalue: lvalue,
						Rvalue: rvalue,
					}
				}
			case MODULE:
				{
					// we dont eat anything we just forward it
					return parser.ParseModule()
				}
			case IMPORT:
				{
					parser.eatToken()
					fileName := parser.CurrentToken()
					parser.eatToken()
					return Import{
						FileName: fileName.Value.(string),
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

				lvalue := parser.ParseAssignment()

				return Assignment{
					Type:   REASSIGNMENT,
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
