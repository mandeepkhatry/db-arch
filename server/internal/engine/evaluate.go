package engine

import (
	"db-arch/server/internal/def"
	"db-arch/server/internal/engine/marshal"
	"db-arch/server/internal/engine/stack"
	"db-arch/server/io"
	"fmt"
	"reflect"
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

var arthmeticExecution = map[string]func(io.Store, string, string, []byte, []byte,
	[]byte, []byte) (roaring.Bitmap, error){

	"=": func(s io.Store, fieldName string, fieldType string, byteOrderedValue []byte,
		dbID []byte, namespaceID []byte, collectionID []byte) (roaring.Bitmap, error) {
		fmt.Println("[[evaluate.go/arthemticExecution]]")
		rb := roaring.New()

		fmt.Println("[[evaluate.go]]/44]", dbID, collectionID, namespaceID)
		indexKey := []byte(def.INDEX_KEY + string(dbID) + ":" + string(collectionID) + ":" + string(namespaceID) + ":" + fieldName + ":" + fieldType + ":" + string(byteOrderedValue))

		uniqueIDBitmapArray, err := s.Get(indexKey)
		fmt.Println("[[evaluate.go/uniqueIDBitmapArray]]", uniqueIDBitmapArray)
		if len(uniqueIDBitmapArray) == 0 || err != nil {
			return roaring.Bitmap{}, err
		}
		err = rb.UnmarshalBinary(uniqueIDBitmapArray)
		fmt.Println("[[evaluate.go/rb]]", rb)

		if err != nil {
			return roaring.Bitmap{}, err
		}
		fmt.Println("[[evaluate.go/rb]],rb")
		return *rb, nil
	},

	">": func(s io.Store, fieldName string, fieldType string, byteOrderedValue []byte,
		dbID []byte, namespaceID []byte, collectionID []byte) (roaring.Bitmap, error) {
		return roaring.Bitmap{}, nil
	},

	"<": func(s io.Store, fieldName string, fieldType string, byteOrderedValue []byte,
		dbID []byte, namespaceID []byte, collectionID []byte) (roaring.Bitmap, error) {

		return roaring.Bitmap{}, nil

	},

	">=": func(s io.Store, fieldName string, fieldType string, byteOrderedValue []byte,
		dbID []byte, namespaceID []byte, collectionID []byte) (roaring.Bitmap, error) {

		return roaring.Bitmap{}, nil

	},

	"<=": func(s io.Store, fieldName string, fieldType string, byteOrderedValue []byte,
		dbID []byte, namespaceID []byte, collectionID []byte) (roaring.Bitmap, error) {

		return roaring.Bitmap{}, nil

	},

	"!=": func(s io.Store, fieldName string, fieldType string, byteOrderedValue []byte,
		dbID []byte, namespaceID []byte, collectionID []byte) (roaring.Bitmap, error) {

		return roaring.Bitmap{}, nil
	},
}

//EvaluatePostFix evaluates postfix expression returns result
func (e *Engine) EvaluatePostFix(s io.Store, px []string, collectionID []byte) (interface{}, error) {
	var tempStack stack.Stack
	for _, v := range px {
		if _, ok := operators[v]; !ok {
			tempStack.Push(v)
		} else {
			exp1 := tempStack.Pop().(string)
			exp2 := tempStack.Pop().(string)
			//TODO: handle here
			rb1, err := e.EvaluateExpression(s, exp1, collectionID)
			if err != nil {
				var tmp interface{}
				return tmp, err
			}
			rb2, err := e.EvaluateExpression(s, exp2, collectionID)
			if err != nil {
				var tmp interface{}
				return tmp, err
			}
			tempStack.Push(execute[v](&rb1, &rb2))
		}
	}
	return tempStack.Pop(), nil
}

//EvaluateExpression takes in expression and returns roaring bitmap as result
func (e *Engine) EvaluateExpression(s io.Store, exp string, collectionID []byte) (roaring.Bitmap, error) {
	/*
		1. Parse expression to find fieldname, operator, fieldvalue, fieldtype
		2. Based on operator, carry out operations
	*/

	//parse fieldname,operator,fieldvalue
	fieldname, operator, fieldvalue := parseExpressionFields(exp)
	//get fieldtype with ordered value
	typeOfData, byteOrderedData := findTypeOfValue(fieldvalue)
	fmt.Println("[[evaluate.go]]typeOfData,byteorderedData:", typeOfData, byteOrderedData)
	fmt.Println("[[evaluate.go]]131]", e.DBID, e.NamespaceID, collectionID)
	rb, err := arthmeticExecution[operator](s, fieldname, typeOfData, byteOrderedData, e.DBID, e.NamespaceID, collectionID)
	if err != nil {
		return roaring.Bitmap{}, err
	}
	return rb, nil
}

func parseExpressionFields(exp string) (string, string, string) {
	re := regexp.MustCompile(`(!=|>=|>|<=|<|=)`)
	operator := re.Find([]byte(exp)) //get first operator that is matched
	strArr := strings.Split(exp, string(operator))
	return strArr[0], string(operator), strings.Trim(strArr[1], "\"")
}

func findTypeOfValue(value string) (string, []byte) {
	fmt.Println("[[evaluate.go]] findTypeOfValue:", reflect.TypeOf(value))
	return "word", marshal.TypeMarshal("string", value)

	//if strings.Contains(value, "'") || strings.Contains(value, "\"") {
	//	fmt.Println("[[evaluate.go]]findtypeofvalue-string:", value)
	//	return "string", marshal.TypeMarshal("string", value)
	//} else if valid.IsInt(value) {
	//	val, err := strconv.Atoi(value)
	//	fmt.Println("[[evaluate.go]]err-int:", err)
	//	fmt.Println("[[evaluate.go]]findtypeofvalue-int:", val)
	//
	//	return "int", marshal.TypeMarshal("int", val)
	//} else if valid.IsFloat(value) {
	//	val, _ := strconv.ParseFloat(value, 64)
	//	fmt.Println("[[evaluate.go]]findtypeofvalue-float:", val)
	//
	//	return "float", marshal.TypeMarshal("float", val)
	//} else if value == "true" || value == "false" {
	//	val, _ := strconv.ParseBool(value)
	//	return "bool", marshal.TypeMarshal("bool", val)
	//} else {
	//	time, _ := time.Parse(time.RFC3339, value)
	//	if time.String() != "0001-01-01 00:00:00 +0000 UTC" {
	//		return "datetime", marshal.TypeMarshal("datetime", time)
	//	}
	//}
	//return "new_data_type", []byte{}
}
