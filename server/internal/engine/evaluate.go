package engine

import (
	"db-arch/server/internal/engine/stack"
	"fmt"
	"regexp"
	"strings"

	"github.com/RoaringBitmap/roaring"
)

var operators = map[string]bool{
	"OR":  true,
	"AND": true,
	"NOT": true,
}

var arithmeticOperators = map[string]bool{
	"=":  true,
	">=": true,
	"<=": true,
	">":  true,
	"<":  true,
	"!=": true,
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

var arthmeticExecution = map[string]func(string, string, []byte) roaring.Bitmap{
	"=": func(fieldName string, typeOfData string, byteOrderedValue []byte) roaring.Bitmap {

	},
	">": func(fieldName string, typeOfData string, byteOrderedValue []byte) roaring.Bitmap {

	},
	"<": func(fieldName string, typeOfData string, byteOrderedValue []byte) roaring.Bitmap {

	},
	">=": func(fieldName string, typeOfData string, byteOrderedValue []byte) roaring.Bitmap {

	},
	"<=": func(fieldName string, typeOfData string, byteOrderedValue []byte) roaring.Bitmap {

	},
	"!=": func(fieldName string, typeOfData string, byteOrderedValue []byte) roaring.Bitmap {

	},
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
			rb1 := EvaluateExpression(exp1)
			rb2 := EvaluateExpression(exp2)
			tempStack.Push(execute[v](rb1, rb2))
		}
	}
	return tempStack.Pop()
}

//EvaluateExpression takes in expression and returns roaring bitmap as result
func EvaluateExpression(exp string) roaring.Bitmap {
	/*
		1. Parse expression to find fieldname, operator, fieldvalue, fieldtype
		2. Based on operator, carry out operations
	*/

	//parse fieldname,operator,fieldvalue
	fieldname, operator, fieldvalue := ParseExpressionFields(exp)
	//get fieldtype with ordered value
	typeOfData, byteOrderedData := findTypeOfValue(fieldvalue) //TODO: implement this

	rb := arthmeticExecution[operator](fieldname, typeOfData, byteOrderedData)
	return rb
}

func ParseExpressionFields(exp string) (string, string, string) {
	re := regexp.MustCompile(`(!=|>=|>|<=|<|=)`)
	operator := re.Find([]byte(exp)) //get first operator that is matched
	strArr := strings.Split(exp, string(operator))
	return strArr[0], string(operator), strings.Trim(strArr[1], '"')
}

func findTypeOfValue(fieldvalue string) (string, []byte) {
	dataType := fmt.Sprintf("%T", fieldvalue)
	//TODO: someone implement this please :D

}
