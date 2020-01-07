package marshal

import (
	"encoding/binary"
	"fmt"
	"math"
	"strings"
	"time"
)

//Integer Range
const (
	INT_RANGE = int(32768)
)

//TypeMarshal takes type of data and value interface as inputs and returns type specific data  byte
func TypeMarshal(typeOfData string, valueInterface interface{}) []byte {
	if typeOfData == "int" {
		return marshalInt(valueInterface)

	} else if typeOfData == "float" {
		return marshalFloat(valueInterface)

	} else if typeOfData == "string" {
		if len(strings.Split(valueInterface.(string), " ")) > 0 {
			return marshalWords(valueInterface)
		}
		return marshalWords(valueInterface)

	} else if typeOfData == "time.Time" {
		return marshalDateTime(valueInterface)

	} else if typeOfData == "bool" {
		return marshalBool(valueInterface)

	}
	fmt.Println("NEW TYPE OF DATA FOUND")
	fmt.Println("DATA : ", valueInterface)
	return []byte("New Type")
}

func marshalInt(valueInterface interface{}) []byte {
	buf := make([]byte, 8)
	//32768 represents range
	//Make sure to observe this range
	numToConvert := int(valueInterface.(float64)) + INT_RANGE
	binary.BigEndian.PutUint32(buf, uint32(numToConvert))
	return buf
}

func marshalFloat(valueInterface interface{}) []byte {
	buf := make([]byte, 8)
	floatNumber := valueInterface.(float64)
	binary.BigEndian.PutUint64(buf, math.Float64bits(floatNumber))
	if floatNumber < 0 {
		for i, _ := range buf {
			buf[i] = ^buf[i]
		}
	} else {
		buf[0] = buf[0] ^ 0x80
	}

	if buf[0] <= 127 {
		for i, _ := range buf {
			buf[i] = ^buf[i]
		}
		// num := binary.BigEndian.Uint64(buf)
		// fmt.Println("decoded", math.Float64frombits(num))
	} else {
		buf[0] = buf[0] ^ 0x80
		// num := binary.BigEndian.Uint64(buf)
		// fmt.Println("decoded", math.Float64frombits(num))
	}

	return buf
}

//TODO Not confirmed
func marshalWord(valueInterface interface{}) []byte {
	return []byte(valueInterface.(string))
}

//TODO Not confirmed
func marshalWords(valueInterface interface{}) []byte {
	return []byte(valueInterface.(string))

}

func marshalBool(valueInterface interface{}) []byte {
	if valueInterface.(bool) == true {
		return []byte("true")
	}
	return []byte("false")
}

func marshalDateTime(valueInterface interface{}) []byte {
	fmt.Println(valueInterface.(time.Time))
	byteKeyTimestamp, _ := valueInterface.(time.Time).MarshalBinary()
	return byteKeyTimestamp

}