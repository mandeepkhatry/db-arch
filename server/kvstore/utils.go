package kvstore

import (
	"bytes"
	"encoding/binary"
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
