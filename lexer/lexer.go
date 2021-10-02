package lexer

import (
	"strings"
	"unicode"
)

func Contains(array []string, str string) bool {
	for _, a := range array {
		if strings.Compare(a, str) == 0 {
			return true
		}
	}

	return false
}

type Lexer struct {
	code            string
	SplitCode       []string
	tokens          []Token
	currentPosition int
	line            int // get the line and start it
	column          int // the code column --> get the start column
}

func InitLexer(code string) *Lexer {
	return &Lexer{
		code:            code,
		currentPosition: 0,
		SplitCode:       strings.Split(code, "\n"),
	}
}

func (lexer *Lexer) CurrentLexeme() rune {
	return rune(lexer.code[lexer.currentPosition])
}

func (lexer *Lexer) peakAhead() rune {
	return rune(lexer.code[lexer.currentPosition+1])
}

func (lexer *Lexer) eatLexeme() {
	lexer.currentPosition++
}

// this appends the token to the system
func (lexer *Lexer) AppendToken(tokenType TokenType, value interface{}, span int) {

	// use a column span
	lexer.tokens = append(lexer.tokens, Token{
		Type:        tokenType,
		Line:        lexer.line,
		Value:       value,
		ColumnStart: lexer.column,
		Span:        span,
	})

	lexer.column += span
}

func isValidVariableName(lexeme rune, continuous bool) bool {
	_start := unicode.IsLetter(lexeme) || lexeme == '_'

	if continuous {
		return _start || unicode.IsDigit(lexeme)
	}

	return _start
}

func (lexer *Lexer) Lex() []Token {
	// var tokens []Token

	for lexer.currentPosition < len(lexer.code) {
		lexeme := lexer.CurrentLexeme()

		// use a regex to match all the letters here
		// or the lexeme is _

		if isValidVariableName(lexeme, false) {
			_currentPosition := lexer.currentPosition

			for ; lexer.currentPosition < len(lexer.code) && isValidVariableName(lexer.CurrentLexeme(), true); lexer.currentPosition++ {

			}

			// what happens is if we have one stuff eg
			_variable := lexer.code[_currentPosition:lexer.currentPosition]
			// lexer.currentPosition -= 1

			tokenType := VARIABLE

			if Contains(KEYWORDS, _variable) {
				tokenType = KEYWORD
			}

			lexer.AppendToken(tokenType, _variable, len(_variable))
			continue
		} else if unicode.IsDigit(lexeme) {
			_start_position_ := lexer.currentPosition

			// we want to get all the number
			// also check for a . indicating a decimal number also monitor the numbers of zero we have

			for ; lexer.currentPosition < len(lexer.code) && (unicode.IsDigit(lexer.CurrentLexeme()) || lexer.CurrentLexeme() == '.'); lexer.currentPosition++ {

			}

			_number_ := lexer.code[_start_position_:lexer.currentPosition]

			// check if the numbers of . is more than one if so just throw an error here
			if strings.Count(_number_, ".") > 1 {
				panic("Too many points")
			}

			// fmt.Println(_number_)

			lexer.currentPosition -= 1

			lexer.AppendToken(NUMBER, _number_, len(_number_))
		} else if lexeme == ';' {
			// tokens = append(tokens, Token{
			// 	Type:  SEMI_COLON,
			// 	Value: ";",
			// })

			lexer.AppendToken(SEMI_COLON, ";", 1)
		} else if lexeme == '!' {
			_current_column_ := lexer.currentPosition
			// _peakAheadLexeme := lexer.peakAhead()

			if lexer.peakAhead() == '=' {
				lexer.eatLexeme()

				// tokens = append(tokens, Token{
				// 	Type:  CONDITION,
				// 	Value: "!=",
				// })

				lexer.AppendToken(CONDITION, "!=", lexer.currentPosition-1-_current_column_)
			}

		} else if lexeme == '=' {
			// check the next if its a another =
			// _current_column_ := lexer.currentPosition
			// return the condition token
			// _peakAheadLexeme := lexer.peakAhead()

			if lexer.peakAhead() == '=' {
				lexer.eatLexeme()

				// tokens = append(tokens, Token{
				// 	Type:  CONDITION,
				// 	Value: "==",
				// })

				lexer.AppendToken(CONDITION, "==", 2)
			} else {
				// tokens = append(tokens, Token{
				// 	Type:  ASSIGN,
				// 	Value: "=",
				// })

				lexer.AppendToken(ASSIGN, "=", 1)
			}
		} else if strings.Contains(`"'`, string(lexeme)) {
			// we do have the starting point go on until we reach the end of the strings
			_initial_position := lexer.currentPosition + 1

			// eat the current lexeme
			lexer.eatLexeme()

			for ; lexer.currentPosition < len(lexer.code) && lexer.CurrentLexeme() != lexeme; lexer.currentPosition++ {
			}

			_string_value_ := lexer.code[_initial_position:lexer.currentPosition]

			lexer.AppendToken(STRING, _string_value_, lexer.currentPosition-_initial_position-2)
			// tokens = append(tokens, Token{
			// 	Type:  STRING,
			// 	Value: _string_value_,
			// })

		} else if Contains([]string{">", "<"}, string(lexeme)) {
			_peakAheadLexeme := lexer.peakAhead()

			if _peakAheadLexeme == '=' {
				lexer.eatLexeme()

				// tokens = append(tokens, Token{
				// 	Type:  CONDITION,
				// 	Value: string(lexeme) + string(_peakAheadLexeme),
				// })

				lexer.AppendToken(CONDITION, string(lexeme)+string(_peakAheadLexeme), 2)
			} else {
				// tokens = append(tokens, Token{
				// 	Type:  CONDITION,
				// 	Value: string(lexeme),
				// })

				lexer.AppendToken(CONDITION, string(lexeme), 1)
			}
		} else if strings.Contains("+-*/%", string(lexeme)) {
			// tokens = append(tokens, Token{
			// 	Type:  OPERATOR,
			// 	Value: string(lexeme),
			// })

			lexer.AppendToken(OPERATOR, string(lexeme), 1)
		} else if lexeme == '.' {
			// tokens = append(tokens, Token{
			// 	Type:  DOT,
			// 	Value: ".",
			// })

			lexer.AppendToken(DOT, ".", 1)
		} else if lexeme == ':' {
			// return the lexeme
			// tokens = append(tokens, Token{
			// 	Type:  COLON,
			// 	Value: ":",
			// })

			lexer.AppendToken(COLON, ":", 1)
		} else if strings.Contains("[]", string(lexeme)) {
			// return the token now
			// tokens = append(tokens, Token{
			// 	Type:  SQUARE_BRACKET,
			// 	Value: string(lexeme),
			// })

			lexer.AppendToken(SQUARE_BRACKET, string(lexeme), 1)
		} else if lexeme == ',' {
			// add the comma tokens for this very reason
			// tokens = append(tokens, Token{
			// 	Type:  COMMA,
			// 	Value: ",",
			// })

			lexer.AppendToken(COMMA, ",", 1)
		} else if lexeme == '\n' {
			lexer.line += 1
			lexer.column = 0
		} else if unicode.IsSpace(lexeme) {
			lexer.column += 1
		} else if lexeme == '#' {
			lexer.eatLexeme()

			// #[ comments ]# for multiline comments

			// if the nextToken is a [
			// eat all the characters until we find ]#

			_start_position_ := lexer.currentPosition

			for ; lexer.currentPosition < len(lexer.code) && lexer.CurrentLexeme() != '\n'; lexer.currentPosition++ {

			}

			_comment_ := lexer.code[_start_position_:lexer.currentPosition]
			lexer.currentPosition -= 1

			// tokens = append(tokens, Token{
			// 	Type:  COMMENT,
			// 	Value: _comment_,
			// })

			lexer.AppendToken(COMMENT, _comment_, lexer.currentPosition-_start_position_)

		} else if strings.Contains("()", string(lexeme)) {
			// tokens = append(tokens, Token{
			// 	Type:  HALF_CIRCLE_BRACKET,
			// 	Value: string(lexeme),
			// })

			lexer.AppendToken(HALF_CIRCLE_BRACKET, string(lexeme), 1)
		} else if strings.Contains("{}", string(lexeme)) {
			// tokens = append(tokens, Token{
			// 	Type:  CURLY_BRACES,
			// 	Value: string(lexeme),
			// })

			lexer.AppendToken(CURLY_BRACES, string(lexeme), 1)
		}

		lexer.currentPosition += 1
	}
	return lexer.tokens
}
