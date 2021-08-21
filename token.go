package main

var (
	KEYWORDS = []string{FUNC, TRUE, FALSE}
)

type TokenType = int

const (
	NUMBER              TokenType = iota + 1
	STRING                        // 2
	VARIABLE                      // 3
	KEYWORD                       // 4
	ASSIGN                        // 5
	OPERATOR                      // 6
	COMMENT                       // 7
	HALF_CIRCLE_BRACKET           // 8
	CURLY_BRACES                  // 9
	COMMA                         //10
)

type Token struct {
	Type  TokenType
	Value interface{}
}
