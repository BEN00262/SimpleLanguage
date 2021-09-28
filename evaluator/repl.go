package evaluator

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	. "github.com/BEN00262/simpleLang/lexer"
	. "github.com/BEN00262/simpleLang/parser"
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
		lexer := InitLexer(text)
		parser := InitParser(lexer.Lex())
		fmt.Println(_print(evaluator.ReplExecute(parser.Parse())))
		printRepl()
	}

	evaluator.TearDownRepl()
	fmt.Println("Bye!")
}
