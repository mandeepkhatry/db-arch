package marshal

import (
	"encoding/binary"
	"fmt"
	"time"
)

//Integer Range
const (
	INT_RANGE = int(32768)
)

//TypeMarshal takes type of data and value interface as inputs and returns type specific data  byte
func TypeMarshal(typeOfData string, valueInterface interface{}) []byte {
	if typeOfData == "int" {
		buf := make([]byte, 8)
		//32768 represents range
		//Make sure to observe this range
		numToConvert := int(valueInterface.(float64)) + INT_RANGE
		binary.BigEndian.PutUint32(buf, uint32(numToConvert))
		return buf

	} else if typeOfData == "float" {
		return []byte("TODO")

	} else if typeOfData == "string" {
		return []byte(valueInterface.(string))

	} else if typeOfData == "timestamp" {
		byteKeyTimestamp, _ := valueInterface.(time.Time).MarshalBinary()
		return byteKeyTimestamp

	} else if typeOfData == "bool" {
		if valueInterface.(bool) == true {
			return []byte("true")
		}
		return []byte("false")

	}
	fmt.Println("NEW TYPE OF DATA FOUND")
	return []byte("New Type")
}
