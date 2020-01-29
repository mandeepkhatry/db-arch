package engine

import (
	"db-arch/server/internal/engine/stack"

	"github.com/RoaringBitmap/roaring"
)

var operators = map[string]bool{
	"OR":  true,
	"AND": true,
	"NOT": true,
}

var arithmeticOperators=map[string]bool{
	"=":true
}

//TODO add corresponding function to each operator ("AND, "OR", "NOT")
var execute = map[string]func(roaring.Bitmap, roaring.Bitmap) roaring.Bitmap{
	"AND": func(rb1, rb2 roaring.Bitmap) roaring.Bitmap {
		return roaring.FastAnd(rb1, rb2)
	},
	"OR": func(rb1, rb2 roaring.Bitmap) roaring.Bitmap {
		return roaring.FastOr(rb1, rb2)
	},
	//TODO: implement NOT
	//"NOT IN": func(rb1, rb2 roaring.Bitmap) roaring.Bitmap {
	//	return rb1.AndNot()
	//},
}

//EvaluatePostFix evaluates postfix expression returns result
func EvaluatePostFix(px []string) interface{} {
	var tempStack stack.Stack
	for _, v := range px {
		if _, ok := operators[v]; !ok {
			tempStack.Push(v)
		} else {
			exp1 := tempStack.Pop().(string)
			exp2 := tempStack.Pop().(string)

			tempStack.Push(execute[v](exp1, exp2))
		}
	}
	return tempStack.Pop()
}

//EvaluateExpression takes in expression and returns roaring bitmap as result
func EvaluateExpression(exp string) (roaring.Bitmap, error) {
	/*
		1. Parse expression to find fieldname, operator, fieldvalue, fieldtype
		2. Based on operator, carry out operations
	*/

}
