package main

import (
	"context"
	"db-arch/pb/document"
	"encoding/json"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

//find if type of data is float64
func findIfFLoat(typeOfData string) bool {
	if typeOfData == "float64" {
		return true
	}
	return false
}

//find if data is integer type
//Note : data from json even in form of integer is represented as float64 type
func checkIfInt(data float64) bool {
	ipart := int64(data)
	decpart := fmt.Sprintf("%.6g", data-float64(ipart))

	if decpart == "0" {
		return true
	}

	return false
}

//returns type of data with keys as data field and value as type
func findTypeOfData(data []byte) map[string]string {

	var dataInterface map[string]interface{}
	err := json.Unmarshal(data, &dataInterface)

	if err != nil {
		panic(err)
	}
	//typeOfData represents a map with key that represents data field and value that represents type of data
	typeOfData := make(map[string]string)

	for k, v := range dataInterface {
		temp := fmt.Sprintf("%T", v)
		//Note : data from json even in form of integer is represented as float64 type
		if findIfFLoat(temp) {
			if checkIfInt(v.(float64)) {
				typeOfData[k] = "int"
			} else {
				typeOfData[k] = "float64"
			}
		} else {
			typeOfData[k] = temp

		}
	}
	return typeOfData
}

//Server struct
type server struct{}

//DocumentTransfer function that recieves request from client and returns response
func (*server) DocumentTransfer(ctx context.Context, req *document.DocumentTransferRequest) (*document.DocumentTransferResponse, error) {
	fmt.Println("-----Got Document as Request-----")

	database := req.GetRequest().GetDatabase()
	collection := req.GetRequest().GetCollection()
	namespace := req.GetRequest().GetNamespace()
	//data in form of bytes
	data := req.GetRequest().GetData()
	/*
		To unmarshal data use this :
			var dataInterface map[string]interface{}
			err := json.Unmarshal(data, &dataInterface)

			if err != nil {
				panic(err)
			}
	*/
	indices := req.GetRequest().GetIndices()

	fmt.Println("Database : ", database)
	fmt.Println("Collection : ", collection)
	fmt.Println("Namespace : ", namespace)
	fmt.Println("Data : ", data)
	fmt.Println("Indices : ", indices)

	//type of data with keys as data field and value as type
	//eg. map[name:string age:int]
	dataType := findTypeOfData(data)
	fmt.Println("Data Type :\n", dataType)

	//Response to client
	res := &document.DocumentTransferResponse{
		Response: "document recieved by server",
	}

	//variables that needs to be fed to database specific functions
	/*
		Data Field		Variable used		Type
		Database		database			string
		Colection		collection			string
		Namespace		namespace			string
		Data			data				[]byte
		Indices			indices				[]string
		DataTyoe		dataType			map[string]string
	*/

	return res, nil
}

func main() {

	fmt.Println("-------------------Starting GRPC server-------------------")
	//just in case server crashes, get detailed log
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	//create a new listener
	lis, err := net.Listen("tcp", "0.0.0.0:50051") //default port for gRPC is 50051

	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	//New GRPC Server
	s := grpc.NewServer()
	document.RegisterDocumentServiceServer(s, &server{})

	reflection.Register(s)

	err = s.Serve(lis)

	if err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}

}
