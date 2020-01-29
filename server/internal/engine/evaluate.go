package engine

import (
	"db-arch/server/internal/def"
	"db-arch/server/internal/engine/stack"
	"db-arch/server/io"
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

//TODO add corresponding function to each operator ("AND, "OR", "NOT")
var execute = map[string]func(*roaring.Bitmap, *roaring.Bitmap) roaring.Bitmap{
	"AND": func(rb1, rb2 *roaring.Bitmap) roaring.Bitmap {
		return *roaring.FastAnd(rb1, rb2)
	},
	"OR": func(rb1, rb2 *roaring.Bitmap) roaring.Bitmap {
		return *roaring.FastOr(rb1, rb2)
	},
	//TODO: implement NOT
	//"NOT IN": func(rb1, rb2 roaring.Bitmap) roaring.Bitmap {
	//	return rb1.AndNot()
	//},
}

var arthmeticExecution = map[string]func(io.Store, string, string, []byte,[]byte,[]byte,[]byte) (roaring.Bitmap,error){
	"=": func(s io.Store, fieldName string, fieldType string, byteOrderedValue []byte,
		dbID []byte,namespaceID []byte,collectionID []byte) (roaring.Bitmap,error) {

		rb := roaring.New()

		indexKey := []byte(def.INDEX_KEY + string(dbID) + ":" + string(collectionID) + ":" + string(namespaceID) + ":" + fieldName + ":" + fieldType + ":" + string(byteOrderedValue)

		uniqueIDBitmapArray,err:=s.Get(indexKey)
		if len(uniqueIDBitmapArray)==0 || err!=nil{
			return roaring.Bitmap{},err
		}


	},
	">": func(s io.Store, fieldName string, fieldType string, byteOrderedValue []byte,
		dbID []byte,namespaceID []byte,collectionID []byte) (roaring.Bitmap,error) {

	},
	"<": func(s io.Store, fieldName string, fieldType string, byteOrderedValue []byte,
		dbID []byte,namespaceID []byte,collectionID []byte) (roaring.Bitmap,error) {

	},
	">=": func(s io.Store, fieldName string,fieldType string, byteOrderedValue []byte,
		dbID []byte,namespaceID []byte,collectionID []byte) (roaring.Bitmap,error) {

	},
	"<=": func(s io.Store, fieldName string, fieldType string, byteOrderedValue []byte,
		dbID []byte,namespaceID []byte,collectionID []byte) (roaring.Bitmap,error) {

	},
	"!=": func(s io.Store, fieldName string, fieldType string, byteOrderedValue []byte,
		dbID []byte,namespaceID []byte,collectionID []byte) (roaring.Bitmap,error) {

	},
}

//EvaluatePostFix evaluates postfix expression returns result
func (e *Engine) EvaluatePostFix(s io.Store, px []string,collectionID []byte)(interface{},error) {
	var tempStack stack.Stack
	for _, v := range px {
		if _, ok := operators[v]; !ok {
			tempStack.Push(v)
		} else {
			exp1 := tempStack.Pop().(string)
			exp2 := tempStack.Pop().(string)
			//TODO: handle here
			rb1,err := e.EvaluateExpression(s,exp1,collectionID)
			rb2,err := e.EvaluateExpression(s,exp2,collectionID)
			tempStack.Push(execute[v](rb1, rb2))
		}
	}
	return tempStack.Pop()
}

//EvaluateExpression takes in expression and returns roaring bitmap as result
func (e *Engine) EvaluateExpression(s io.Store,exp string,collectionID []byte) (roaring.Bitmap,error) {
	/*
		1. Parse expression to find fieldname, operator, fieldvalue, fieldtype
		2. Based on operator, carry out operations
	*/

	//parse fieldname,operator,fieldvalue
	fieldname, operator, fieldvalue := parseExpressionFields(exp)
	//get fieldtype with ordered value
	typeOfData, byteOrderedData := findTypeOfValue(fieldvalue) //TODO: implement this

	rb,err := arthmeticExecution[operator](s,fieldname, typeOfData, byteOrderedData,e.DBID,e.NamespaceIdentifier,collectionID)
	if err!=nil{
		return roaring.Bitmap{},err
	}
	return rb,nil
}

func parseExpressionFields(exp string) (string, string, string) {
	re := regexp.MustCompile(`(!=|>=|>|<=|<|=)`)
	operator := re.Find([]byte(exp)) //get first operator that is matched
	strArr := strings.Split(exp, string(operator))
	return strArr[0], string(operator), strings.Trim(strArr[1], "\"")
}

func findTypeOfValue(fieldvalue string) (string, []byte) {
	dataType := fmt.Sprintf("%T", fieldvalue)
	//TODO: someone implement this please :D

}
