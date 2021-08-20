package main

type Evaluator struct {
	program      *ProgramNode
	symbolsTable *SymbolsTable
}

func initEvaluator(program *ProgramNode) *Evaluator {
	return &Evaluator{
		program:      program,
		symbolsTable: initSymbolsTable(),
	}
}

func (eval *Evaluator) walkTree(node interface{}) error {
	switch node.(type) {
	case VariableNode:
		{
			// get the value from the symbols Table
		}
	case NumberNode:
		{
			// return the value
		}
	}

	return nil
}

func (eval *Evaluator) Evaluate() interface{} {
	eval.symbolsTable.pushContext()

	for _, node := range eval.program.Nodes {
		err := eval.walkTree(node)

		if err != nil {
			panic(err.Error())
		}
	}

	eval.symbolsTable.popContext()
	return nil
}
