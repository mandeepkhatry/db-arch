package kvstore

import (
	"bytes"
	"db-arch/server/kvstore/marshal"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"math/rand"
	"net"
	"os"
	"time"
)

/*
	All utility functions are defined here!
*/

//getApplicationSpecificType return application specific data type
func getApplicationSpecificType(typeOfData string, valueInterface interface{}) string {
	if typeOfData == "string" {
		//TODO seperate for word and words
		return "word"
	} else if typeOfData == "float64" || typeOfData == "float32" {
		return "double"
	} else if typeOfData == "bool" {
		return "bool"
	} else if typeOfData == "time.Time" {
		return "datatime"
	}

	//New type
	return "new type"
}

//getMACAddress return 3 byte MAC address of current machine
func getMACAddress() []byte {
	interfaces, err := net.Interfaces()
	var addr string
	if err != nil {
		panic(err)
	}

	for _, i := range interfaces {
		if i.Flags&net.FlagUp != 0 && bytes.Compare(i.HardwareAddr, nil) != 0 {
			addr = i.HardwareAddr.String()
			break
		}
	}
	return []byte(addr)
}

//getUnixTimeStamp returns 4 byte UNIX timestamp
func getUnixTimestamp() []byte {
	currentTimestamp := time.Now().UnixNano()
	timeInBytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(timeInBytes, uint32(currentTimestamp))
	return timeInBytes
}

//TODO: generate multiple counters at a time for batch entries
//generateRandomCount generates a 4 bytes random integer
func generateRandomCount() []byte {
	/*
		Generate a random 32bit uint value
		(0 to 4294967295)
	*/
	rand.Seed(time.Now().UnixNano())
	count := rand.Uint32()
	countInBytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(countInBytes, uint32(count))
	return countInBytes
}

//getProcessID returns 2 byte current processID
func getProcessID() []byte {
	processIDInBytes := make([]byte, 2)
	binary.LittleEndian.PutUint16(processIDInBytes, uint16(os.Getpid()))
	return processIDInBytes
}

//TODO: make this func generic
//generateKey
func generateKey(dbID []byte, collectionID []byte, namespaceID []byte, uniqueID []byte) []byte {
	//key := ""
	//key = string(dbID) + ":" + string(collectionID) + ":" + string(namespaceID) + ":" + string(uniqueID)
	//return []byte(key)
	key := append(dbID, []byte(":")...)
	key = append(key, collectionID...)
	key = append(key, []byte(":")...)
	key = append(key, namespaceID...)
	key = append(key, []byte(":")...)
	key = append(key, uniqueID...)
	key = append(key, []byte(":")...)
	return key
}

//findIfFloat finds if type of data is float64
func findIfFLoat(typeOfData string) bool {
	if typeOfData == "float64" {
		return true
	}
	return false
}

//checkIfInt finds if data is integer type
//Note : data from json even in form of integer is represented as float64 type
func checkIfInt(data float64) bool {
	ipart := int64(data)
	decpart := fmt.Sprintf("%.6g", data-float64(ipart))

	if decpart == "0" {
		return true
	}

	return false
}

//FindTypeOfData returns type of data with keys as data field and value as type and type specific data bytes
func findTypeOfData(data map[string][]byte) (map[string]string, map[string][]byte) {

	//typeOfData represents a map with key that represents data field and value that represents type of data
	typeOfData := make(map[string]string)
	var valueInterface interface{}

	newData := make(map[string][]byte)

	for k, v := range data {
		err := json.Unmarshal(v, &valueInterface)
		if err != nil {
			panic(err)
		}

		dataType := fmt.Sprintf("%T", valueInterface)
		//Note : data from json even in form of integer is represented as float64 type
		if findIfFLoat(dataType) {
			if checkIfInt(valueInterface.(float64)) {
				typeOfData[k] = getApplicationSpecificType("int", valueInterface)
				newData[k] = marshal.TypeMarshal("int", valueInterface)
			} else {
				typeOfData[k] = getApplicationSpecificType("infloat64t", valueInterface)
				newData[k] = marshal.TypeMarshal("float", valueInterface)
			}

		} else if dataType == "string" {
			time, _ := time.Parse(time.RFC3339, valueInterface.(string))

			if time.String() == "0001-01-01 00:00:00 +0000 UTC" {
				stringType := getApplicationSpecificType(dataType, valueInterface)
				typeOfData[k] = stringType
				newData[k] = marshal.TypeMarshal(stringType, valueInterface)
			} else {
				timeType := getApplicationSpecificType(fmt.Sprintf("%T", time), valueInterface)
				typeOfData[k] = timeType
				newData[k] = marshal.TypeMarshal(timeType, valueInterface)
			}
		} else {
			newData[k] = marshal.TypeMarshal(dataType, valueInterface)
			typeOfData[k] = getApplicationSpecificType(dataType, valueInterface)

		}
	}
	return typeOfData, newData
}
