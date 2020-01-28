package postfix

import (
	"db-arch/server/internal/engine/stack"
)

var precedence = map[string]int{
	"OR":  1,
	"AND": 2,
	"NOT": 3,
}

//FindWeight returns weight of operators
func FindWeight(s string) int {
	return precedence[s]
}

//FindHigherPrecedence return true if s1 has higher precedence over s2
func FindHigherPrecedence(s1 string, s2 string) bool {
	return FindWeight(s1) >= FindWeight(s2)
}

//InfixToPostfix returns postfix expressions in form of array of string
func InfixToPostfix(s []string) []string {
	var tempStack stack.Stack
	postfix := make([]string, 0)

	for _, v := range s {
		if v == "(" {
			tempStack.Push(v)
		} else if v == ")" {
			for !tempStack.Empty() && tempStack.Top().(string) != "(" {
				top := tempStack.Top().(string)
				tempStack.Pop()
				postfix = append(postfix, top)
			}
			if tempStack.Top() == "(" {
				tempStack.Pop()
			}
		} else if _, ok := precedence[v]; ok {
			//code block for operators "AND", "OR", "NOT"
			for !tempStack.Empty() && FindHigherPrecedence(tempStack.Top().(string), v) {
				top := tempStack.Top().(string)
				tempStack.Pop()
				postfix = append(postfix, top)
			}
			tempStack.Push(v)
		} else {
			//code block for operands
			postfix = append(postfix, v)
		}
	}
	//remaining elements from stack
	for !tempStack.Empty() {
		top := tempStack.Top().(string)
		tempStack.Pop()
		postfix = append(postfix, top)
	}
	return postfix
}
