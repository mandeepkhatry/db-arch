package engine

import (
	"db-arch/server/internal/engine/stack"
)

var operators = map[string]bool{
	"OR":  true,
	"AND": true,
	"NOT": true,
	"(":   true,
	")":   true,
}

//TODO add corresponding function to each operator ("AND, "OR", "NOT")
var execute = map[string]func(string, string) string{
	"AND": func(op1, op2 string) string {
		op1result := OpeartorBasedSearch(field, operator, value)
		return op1 + op2 + " anding "
	},
	"OR": func(op1, op2 string) string {
		return op1 + op2 + " oring "
	},
	"NOT": func(op1, op2 string) string {
		return op1 + op2 + " not "
	},
}

//EvaluatePostFix evaluates postfix expression returns result
func EvaluatePostFix(px []string) interface{} {
	var tempStack stack.Stack
	for _, v := range px {
		if _, ok := operators[v]; !ok {
			tempStack.Push(v)
		} else {
			op1 := tempStack.Pop().(string)
			op2 := tempStack.Pop().(string)
			tempStack.Push(execute[v](op1, op2))
		}
	}
	return tempStack.Pop()
}
