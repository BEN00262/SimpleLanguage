package lexer

import (
	"strconv"
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
	currentPosition int
}

func InitLexer(code string) *Lexer {
	return &Lexer{
		code:            code,
		currentPosition: 0,
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

func (lexer *Lexer) Lex() []Token {
	var tokens []Token

	for lexer.currentPosition < len(lexer.code) {
		lexeme := lexer.CurrentLexeme()

		// use a regex to match all the letters here

		if unicode.IsLetter(lexeme) {
			_currentPosition := lexer.currentPosition
			_currentLexeme := lexer.CurrentLexeme()

			lexer.currentPosition += 1 // eat the current lexeme

			for ; lexer.currentPosition < len(lexer.code) && unicode.IsLetter(_currentLexeme); lexer.currentPosition++ {
				_currentLexeme = lexer.CurrentLexeme()
			}

			lexer.currentPosition -= 1

			_variable := lexer.code[_currentPosition:lexer.currentPosition]
			tokenType := VARIABLE

			if Contains(KEYWORDS, _variable) {
				tokenType = KEYWORD
			}

			tokens = append(tokens, Token{
				Type:  tokenType,
				Value: _variable,
			})
			continue
		} else if unicode.IsDigit(lexeme) {
			_start_position_ := lexer.currentPosition

			for ; lexer.currentPosition < len(lexer.code) && unicode.IsDigit(lexer.CurrentLexeme()); lexer.currentPosition++ {

			}

			_number_ := lexer.code[_start_position_:lexer.currentPosition]

			lexer.currentPosition -= 1

			intValue, err := strconv.Atoi(_number_)

			if err != nil {
				panic(err.Error())
			}

			tokens = append(tokens, Token{
				Type:  NUMBER,
				Value: intValue,
			})
		} else if lexeme == ';' {
			tokens = append(tokens, Token{
				Type:  SEMI_COLON,
				Value: ";",
			})
		} else if lexeme == '!' {
			_peakAheadLexeme := lexer.peakAhead()

			if _peakAheadLexeme == '=' {
				lexer.eatLexeme()

				tokens = append(tokens, Token{
					Type:  CONDITION,
					Value: "!=",
				})
			}

		} else if lexeme == '=' {
			// check the next if its a another =
			// return the condition token
			_peakAheadLexeme := lexer.peakAhead()

			if _peakAheadLexeme == '=' {
				lexer.eatLexeme()

				tokens = append(tokens, Token{
					Type:  CONDITION,
					Value: "==",
				})
			} else {
				tokens = append(tokens, Token{
					Type:  ASSIGN,
					Value: "=",
				})
			}
		} else if strings.Contains(`"'`, string(lexeme)) {
			// we do have the starting point go on until we reach the end of the strings
			_initial_position := lexer.currentPosition + 1

			// eat the current lexeme
			lexer.eatLexeme()

			for ; lexer.currentPosition < len(lexer.code) && lexer.CurrentLexeme() != lexeme; lexer.currentPosition++ {
			}

			_string_value_ := lexer.code[_initial_position:lexer.currentPosition]

			tokens = append(tokens, Token{
				Type:  STRING,
				Value: _string_value_,
			})

		} else if Contains([]string{">", "<"}, string(lexeme)) {
			// return a conditional here
			// match further down the line other stuff
			_peakAheadLexeme := lexer.peakAhead()

			if _peakAheadLexeme == '=' {
				lexer.eatLexeme()

				tokens = append(tokens, Token{
					Type:  CONDITION,
					Value: string(lexeme) + string(_peakAheadLexeme),
				})
			} else {
				tokens = append(tokens, Token{
					Type:  CONDITION,
					Value: string(lexeme),
				})
			}
		} else if strings.Contains("+-*/", string(lexeme)) {
			tokens = append(tokens, Token{
				Type:  OPERATOR,
				Value: string(lexeme),
			})
		} else if lexeme == ':' {
			// return the lexeme
			tokens = append(tokens, Token{
				Type:  COLON,
				Value: ":",
			})
		} else if strings.Contains("[]", string(lexeme)) {
			// return the token now
			tokens = append(tokens, Token{
				Type:  SQUARE_BRACKET,
				Value: string(lexeme),
			})
		} else if lexeme == ',' {
			// add the comma tokens for this very reason
			tokens = append(tokens, Token{
				Type:  COMMA,
				Value: ",",
			})
		} else if lexeme == '#' {
			lexer.eatLexeme()

			// #[ comments ]# for multiline comments

			// if the nextToken is a [
			// eat all the characters until we find ]#

			_currentLexeme := lexer.CurrentLexeme()
			_currentPosition := lexer.currentPosition

			for ; _currentLexeme != '\n'; lexer.currentPosition++ {
				_currentLexeme = lexer.CurrentLexeme()
			}

			lexer.currentPosition -= 1

			tokens = append(tokens, Token{
				Type:  COMMENT,
				Value: lexer.code[_currentPosition:lexer.currentPosition],
			})
			continue
		} else if strings.Contains("()", string(lexeme)) {
			tokens = append(tokens, Token{
				Type:  HALF_CIRCLE_BRACKET,
				Value: string(lexeme),
			})
		} else if strings.Contains("{}", string(lexeme)) {
			tokens = append(tokens, Token{
				Type:  CURLY_BRACES,
				Value: string(lexeme),
			})
		}

		lexer.currentPosition += 1
	}

	return tokens
}
