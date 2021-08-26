package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func get(r *bufio.Reader) string {
	t, _ := r.ReadString('\n')
	return strings.TrimSpace(t)
}

func shouldContinue(text string) bool {
	if strings.EqualFold("exit", text) {
		return false
	}
	return true
}

func printRepl() {
	fmt.Print("go-repl> ")
}

func REPL() {
	evaluator := NewEvaluatorContext()

	// inject all global functions
	evaluator.InitGlobalScope()

	LoadGlobalsToContext(evaluator)

	reader := bufio.NewReader(os.Stdin)

	printRepl()
	text := get(reader)
	for ; shouldContinue(text); text = get(reader) {
		lexer := initLexer(text)
		parser := initParser(lexer.Lex())
		fmt.Println(evaluator.ReplExecute(parser.Parse()))
		printRepl()
	}

	evaluator.TearDownRepl()
	fmt.Println("Bye!")
}
