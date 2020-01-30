package engine

import (
	"bytes"
	"db-arch/server/internal/def"
	marshal2 "db-arch/server/internal/engine/marshal"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"math/rand"
	"net"
	"os"
	"strconv"
	"strings"
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

	//replace MAC address characters with number
	addr = strings.Replace(addr, "A", "2", -1)
	addr = strings.Replace(addr, "B", "21", -1)
	addr = strings.Replace(addr, "C", "3", -1)
	addr = strings.Replace(addr, "D", "31", -1)
	addr = strings.Replace(addr, "E", "4", -1)
	addr = strings.Replace(addr, "F", "41", -1)
	addr = strings.Replace(addr, ":", "", -1)

	addrInt, err := strconv.Atoi(addr)
	if err != nil {
		panic(err)
	}

	addrInBytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(addrInBytes, uint32(addrInt))
	return addrInBytes
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
//func generateKey(dbID []byte, collectionID []byte, namespaceID []byte, uniqueID []byte) []byte {
//	//key := ""
//	//key = string(dbID) + ":" + string(collectionID) + ":" + string(namespaceID) + ":" + string(uniqueID)
//	//return []byte(key)
//	key:=append(dbID,[]byte(":")...)
//	key=append(key,collectionID...)
//	key=append(key,[]byte(":")...)
//	key=append(key,namespaceID...)
//	key=append(key,[]byte(":")...)
//	key=append(key,uniqueID...)
//	return key
//}

func generateKey(args ...[]byte) []byte {
	key := ""
	length := len(args)
	for i := 0; i < length; i++ {
		key += string(args[i])
		if i < (length - 1) {
			key += string(":")
		} else {
			break
		}
	}
	return []byte(key)
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
		if findIfFLoat(dataType) == true {
			if checkIfInt(valueInterface.(float64)) {
				typeOfData[k] = def.ApplicationSpecificType["int"]
				newData[k] = marshal2.TypeMarshal("int", valueInterface)
			} else {
				typeOfData[k] = def.ApplicationSpecificType["float64"]
				newData[k] = marshal2.TypeMarshal("float64", valueInterface)
			}

		} else if dataType == "string" {
			time, _ := time.Parse(time.RFC3339, valueInterface.(string))

			if time.String() == "0001-01-01 00:00:00 +0000 UTC" {
				typeOfData[k] = def.ApplicationSpecificType["string"]
				newData[k] = marshal2.TypeMarshal("string", valueInterface)
			} else {
				var timeInterface interface{}
				timeInterface = time
				timeType := def.ApplicationSpecificType[fmt.Sprintf("%T", time)]
				typeOfData[k] = timeType
				newData[k] = marshal2.TypeMarshal("time.Time", timeInterface)
			}
		} else {
			newData[k] = marshal2.TypeMarshal(dataType, valueInterface)
			typeOfData[k] = def.ApplicationSpecificType[dataType]

		}
	}
	return typeOfData, newData
}
