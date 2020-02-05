package engine

import (
	"bytes"
	"db-arch/server/internal/def"
	"db-arch/server/internal/engine/formatter"
	"db-arch/server/internal/engine/marshal"
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
var execute = map[string]func(roaring.Bitmap, roaring.Bitmap) roaring.Bitmap{
	"AND": func(rb1, rb2 roaring.Bitmap) roaring.Bitmap {
		return *roaring.FastAnd(&rb1, &rb2)
	},
	"OR": func(rb1, rb2 roaring.Bitmap) roaring.Bitmap {

		return *roaring.FastOr(&rb1, &rb2)
	},
	//TODO: implement NOT
	//"NOT IN": func(rb1, rb2 roaring.Bitmap) roaring.Bitmap {
	//	return rb1.AndNot()
	//},
}

var airthmeticExecution = map[string]func(io.Store, string, string, []byte, []byte,
	[]byte, []byte) (roaring.Bitmap, error){

	"=": func(s io.Store, fieldName string, fieldType string, byteOrderedValue []byte,
		dbID []byte, namespaceID []byte, collectionID []byte) (roaring.Bitmap, error) {
		fmt.Println("[[evaluate.go/airthemeticExecution=]]")
		rb := roaring.New()

		fmt.Println("[[evaluate.go]]/44]", dbID, collectionID, namespaceID)
		indexKey := []byte(def.INDEX_KEY + string(dbID) + ":" + string(collectionID) + ":" + string(namespaceID) + ":" + fieldName + ":" + fieldType + ":" + string(byteOrderedValue))
		fmt.Println("INDEX KEY IS ", string(indexKey))
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

	//TODO: discuss memory related issue here
	">": func(s io.Store, fieldName string, fieldType string, byteOrderedValue []byte,
		dbID []byte, namespaceID []byte, collectionID []byte) (roaring.Bitmap, error) {
		fmt.Println("[[evaluate.go/airthmeticExecution>]]")
		rb := roaring.New()

		startKey := []byte(def.INDEX_KEY + string(dbID) + ":" + string(collectionID) + ":" + string(namespaceID) + ":" + fieldName + ":" + fieldType + ":" + string(byteOrderedValue))
		prefix := []byte(def.INDEX_KEY + string(dbID) + ":" + string(collectionID) + ":" + string(namespaceID) + ":" + fieldName + ":" + fieldType + ":")

		keys, values, err := s.PrefixScan(startKey, prefix, 0)
		fmt.Println("KEYS : ", keys)

		if len(values) == 0 || err != nil {
			return roaring.Bitmap{}, err
		}

		if bytes.Compare(keys[0], startKey) == 0 {
			values = values[1:]
		}

		for _, v := range values {
			if len(rb.ToArray()) == 0 {
				err = rb.UnmarshalBinary(v)
				if err != nil {
					return roaring.Bitmap{}, err
				}
			}
			tempRb := roaring.New()
			err = tempRb.UnmarshalBinary(v)
			if err != nil {
				return roaring.Bitmap{}, err
			}
			fmt.Println("TEMP RB is ", tempRb)
			rb = roaring.FastOr(rb, tempRb)
		}

		fmt.Println("[[evaluate.go/rb]]", rb)

		return *rb, err

	},

	"<": func(s io.Store, fieldName string, fieldType string, byteOrderedValue []byte,
		dbID []byte, namespaceID []byte, collectionID []byte) (roaring.Bitmap, error) {

		fmt.Println("[[evaluate.go/airthmeticExecution<]]")
		rb := roaring.New()

		endKey := []byte(def.INDEX_KEY + string(dbID) + ":" + string(collectionID) + ":" + string(namespaceID) + ":" + fieldName + ":" + fieldType + ":" + string(byteOrderedValue))
		prefix := []byte(def.INDEX_KEY + string(dbID) + ":" + string(collectionID) + ":" + string(namespaceID) + ":" + fieldName + ":" + fieldType + ":")

		keys, values, err := s.ReversePrefixScan(endKey, prefix, 0)

		if len(values) == 0 || err != nil {
			return roaring.Bitmap{}, err
		}

		if bytes.Compare(keys[0], endKey) == 0 {
			values = values[1:]
		}

		for _, v := range values {
			if len(rb.ToArray()) == 0 {
				err = rb.UnmarshalBinary(v)
				if err != nil {
					return roaring.Bitmap{}, err
				}
			}
			tempRb := roaring.New()
			err = tempRb.UnmarshalBinary(v)
			if err != nil {
				return roaring.Bitmap{}, err
			}
			fmt.Println("TEMP RB is ", tempRb)
			rb = roaring.FastOr(rb, tempRb)
		}

		fmt.Println("[[evaluate.go/rb]]", rb)

		return *rb, err

	},

	">=": func(s io.Store, fieldName string, fieldType string, byteOrderedValue []byte,
		dbID []byte, namespaceID []byte, collectionID []byte) (roaring.Bitmap, error) {
		fmt.Println("[[evaluate.go/airthmeticExecution>=]]")
		rb := roaring.New()

		startKey := []byte(def.INDEX_KEY + string(dbID) + ":" + string(collectionID) + ":" + string(namespaceID) + ":" + fieldName + ":" + fieldType + ":" + string(byteOrderedValue))
		prefix := []byte(def.INDEX_KEY + string(dbID) + ":" + string(collectionID) + ":" + string(namespaceID) + ":" + fieldName + ":" + fieldType + ":")

		fmt.Println("STARTKEY", string(startKey))
		fmt.Println("PREFIX", string(prefix))
		_, uniqueIDBitmapArray, err := s.PrefixScan(startKey, prefix, 0)
		fmt.Println("UNIQUEID BITMAP ARRAY is ", uniqueIDBitmapArray)
		if len(uniqueIDBitmapArray) == 0 || err != nil {
			return roaring.Bitmap{}, err
		}

		for _, v := range uniqueIDBitmapArray {
			if len(rb.ToArray()) == 0 {
				err = rb.UnmarshalBinary(v)
				if err != nil {
					return roaring.Bitmap{}, err
				}
			}
			tempRb := roaring.New()
			err = tempRb.UnmarshalBinary(v)
			if err != nil {
				return roaring.Bitmap{}, err
			}
			fmt.Println("TEMP RB is ", tempRb)
			rb = roaring.FastOr(rb, tempRb)
		}

		fmt.Println("[[evaluate.go/rb]]", rb)

		return *rb, err

	},

	"<=": func(s io.Store, fieldName string, fieldType string, byteOrderedValue []byte,
		dbID []byte, namespaceID []byte, collectionID []byte) (roaring.Bitmap, error) {

		fmt.Println("[[evaluate.go/airthmeticExecution<]]")
		rb := roaring.New()

		endKey := []byte(def.INDEX_KEY + string(dbID) + ":" + string(collectionID) + ":" + string(namespaceID) + ":" + fieldName + ":" + fieldType + ":" + string(byteOrderedValue))
		prefix := []byte(def.INDEX_KEY + string(dbID) + ":" + string(collectionID) + ":" + string(namespaceID) + ":" + fieldName + ":" + fieldType + ":")

		_, uniqueIDBitmapArray, err := s.ReversePrefixScan(endKey, prefix, 0)

		if len(uniqueIDBitmapArray) == 0 || err != nil {
			return roaring.Bitmap{}, err
		}

		for _, v := range uniqueIDBitmapArray {
			if len(rb.ToArray()) == 0 {
				err = rb.UnmarshalBinary(v)
				if err != nil {
					return roaring.Bitmap{}, err
				}
			}
			tempRb := roaring.New()
			err = tempRb.UnmarshalBinary(v)
			if err != nil {
				return roaring.Bitmap{}, err
			}
			fmt.Println("TEMP RB is ", tempRb)
			rb = roaring.FastOr(rb, tempRb)
		}

		fmt.Println("[[evaluate.go/rb]]", rb)

		return *rb, err

	},

	"!=": func(s io.Store, fieldName string, fieldType string, byteOrderedValue []byte,
		dbID []byte, namespaceID []byte, collectionID []byte) (roaring.Bitmap, error) {

		fmt.Println("[[evaluate.go/airthmeticExecution<]]")
		rb := roaring.New()

		endKey := []byte(def.INDEX_KEY + string(dbID) + ":" + string(collectionID) + ":" + string(namespaceID) + ":" + fieldName + ":" + fieldType + ":" + string(byteOrderedValue))
		prefix := []byte(def.INDEX_KEY + string(dbID) + ":" + string(collectionID) + ":" + string(namespaceID) + ":" + fieldName + ":" + fieldType + ":")

		keysleft, uniqueIDBitmapArray, err := s.ReversePrefixScan(endKey, prefix, 0)

		if len(uniqueIDBitmapArray) == 0 || err != nil {
			return roaring.Bitmap{}, err
		}

		if bytes.Compare(keysleft[0], endKey) == 0 {
			uniqueIDBitmapArray = uniqueIDBitmapArray[1:]
		}

		keysRight, valuesFor, err := s.PrefixScan(endKey, prefix, 0)

		if len(valuesFor) == 0 || err != nil {
			return roaring.Bitmap{}, err
		}

		if bytes.Compare(keysRight[0], endKey) == 0 {
			valuesFor = valuesFor[1:]
		}

		uniqueIDBitmapArray = append(uniqueIDBitmapArray, valuesFor...)

		for _, v := range uniqueIDBitmapArray {
			if len(rb.ToArray()) == 0 {
				err = rb.UnmarshalBinary(v)
				if err != nil {
					return roaring.Bitmap{}, err
				}
			}
			tempRb := roaring.New()
			err = tempRb.UnmarshalBinary(v)
			if err != nil {
				return roaring.Bitmap{}, err
			}
			fmt.Println("TEMP RB is ", tempRb)
			rb = roaring.FastOr(rb, tempRb)
		}

		fmt.Println("[[evaluate.go/rb]]", rb)

		return *rb, err

	},
}

//EvaluatePostFix evaluates postfix expression returns result
func (e *Engine) EvaluatePostFix(s io.Store, px []string, collectionID []byte) (interface{}, error) {
	if len(px) == 1 {
		rb, err := e.EvaluateExpression(s, px[0], collectionID)
		if err != nil {
			var tmp interface{}
			return tmp, err
		}
		var result interface{}
		result = rb
		return result, nil
	}
	var tempStack stack.Stack
	for _, v := range px {
		if _, ok := operators[v]; !ok {
			tempStack.Push(v)
		} else {
			exp1 := tempStack.Pop()
			exp2 := tempStack.Pop()
			exp1Type := fmt.Sprintf("%T", exp1)
			fmt.Println("Data type : ", exp1Type)
			exp2Type := fmt.Sprintf("%T", exp2)
			fmt.Println("Data type : ", exp2Type)

			var rb1 roaring.Bitmap
			var rb2 roaring.Bitmap
			var err error

			if exp1Type == "string" {
				//TODO: handle here
				rb1, err = e.EvaluateExpression(s, exp1.(string), collectionID)
				if err != nil {
					var tmp interface{}
					return tmp, err
				}
			} else {
				rb1 = exp1.(roaring.Bitmap)
			}

			if exp2Type == "string" {
				rb2, err = e.EvaluateExpression(s, exp2.(string), collectionID)
				if err != nil {
					var tmp interface{}
					return tmp, err
				}

			} else {
				rb2 = exp2.(roaring.Bitmap)
			}

			tempStack.Push(execute[v](rb1, rb2))
		}
	}

	fmt.Println("[[evaluate.go]] EVALUATE POSTFIX")
	result := tempStack.Pop()

	return result, nil

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

	fmt.Println("Operator is ", operator)
	rb, err := airthmeticExecution[operator](s, fieldname, typeOfData, byteOrderedData, e.DBID, e.NamespaceID, collectionID)

	if err != nil {
		return roaring.Bitmap{}, err
	}

	return rb, nil
}

func parseExpressionFields(exp string) (string, string, string) {
	re := regexp.MustCompile(`(!=|>=|>|<=|<|=)`)
	operator := re.Find([]byte(exp)) //get first operator that is matched
	strArr := strings.Split(exp, string(operator))
	return strArr[0], string(operator), (strArr[1])
}

func findTypeOfValue(value string) (string, []byte) {
	fmt.Println("VALUE is ", value)
	datatype, formattedData, err := formatter.FormatData(value)
	fmt.Println("DATATYPE : ", datatype)
	fmt.Println("FORMATTED DATE : ", formattedData)
	if err != nil {
		panic(err)
	}

	specificDataType := def.ApplicationSpecificType[datatype]
	fmt.Println("check")

	return specificDataType, marshal.TypeMarshal(datatype, formattedData)

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
