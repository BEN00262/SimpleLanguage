package evaluator

import (
	. "github.com/BEN00262/simpleLang/lexer"
	. "github.com/BEN00262/simpleLang/parser"
)

// cache the codeString thing
// a simple one line cache
type CodeCache struct {
	CachedCodeString string
	Program          *ProgramNode
	Lexed            *Lexer
}

func (cache *CodeCache) IsFresh(code_string string) bool {
	return cache.CachedCodeString == code_string
}

func (cache *CodeCache) UpdateCache(code_string string, lexed *Lexer, program *ProgramNode) {
	cache.CachedCodeString = code_string
	cache.Program = program
	cache.Lexed = lexed
}
