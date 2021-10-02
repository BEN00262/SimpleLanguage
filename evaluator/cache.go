package evaluator

import (
	Lexer "github.com/BEN00262/simpleLang/lexer"
	Parser "github.com/BEN00262/simpleLang/parser"
)

// cache the codeString thing
// a simple one line cache
type CodeCache struct {
	CachedCodeString string
	Program          *Parser.ProgramNode
	Lexed            *Lexer.Lexer
}

func (cache *CodeCache) IsFresh(code_string string) bool {
	return cache.CachedCodeString == code_string
}

func (cache *CodeCache) UpdateCache(code_string string, lexed *Lexer.Lexer, program *Parser.ProgramNode) {
	cache.CachedCodeString = code_string
	cache.Program = program
	cache.Lexed = lexed
}
