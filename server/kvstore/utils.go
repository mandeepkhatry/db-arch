package kvstore

import (
	"bytes"
	"db-arch/server/kvstore/marshal"
	"encoding/binary"
	"encoding/json"
	"math/rand"
	"net"
	"os"
	"time"
)

/*
	All utility functions are defined here!
*/

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
	key:=append(dbID,[]byte(":")...)
	key=append(key,collectionID...)
	key=append(key,[]byte(":")...)
	key=append(key,namespaceID...)
	key=append(key,[]byte(":")...)
	key=append(key,uniqueID...)
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


//FindTypeOfData returns type of data with keys as data field and value as type
func findTypeOfData(data map[string][]byte) (map[string]string, map[string][]byte) {

	//typeOfData represents a map with key that represents data field and value that represents type of data
	typeOfData := make(map[string]string)
	var tempInterface interface{}

	new_data := make(map[string][]byte)

	for k, v := range data {
		err := json.Unmarshal(v, &tempInterface)
		if err != nil {
			panic(err)
		}

		dataType := fmt.Sprintf("%T", tempInterface)
		//Note : data from json even in form of integer is represented as float64 type
		if findIfFLoat(dataType) {
			if checkIfInt(tempInterface.(float64)) {
				typeOfData[k] = "int"
				new_data[k] = marshal.TypeMarshal("int", tempInterface)
			} else {
				typeOfData[k] = "float64"
				new_data[k] = marshal.TypeMarshal("float64", tempInterface)
			}
		} else {
			new_data[k] = marshal.TypeMarshal(dataType, tempInterface)
			typeOfData[k] = dataType

		}
	}
	return typeOfData, new_data
}
