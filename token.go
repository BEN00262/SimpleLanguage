package main

type TokenType = int

const (
	NUMBER TokenType = iota + 1
	STRING
	VARIABLE
	KEYWORDS
)

type Token struct {
	Type  TokenType
	Value interface{}
}
