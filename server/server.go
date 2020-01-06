package main

import (
	"context"
	"db-arch/pb/document"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"

	"db-arch/server/kvstore"
	"db-arch/server/marshal"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var store kvstore.StoreClient

//TODO Business Logic to map golang specific type

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
	fmt.Println("Data (DATBASE PURPOSE): ", data)
	fmt.Println("Indices : ", indices)

	//type of data with keys as data field and value as type
	//eg. map[name:string age:int]
	dataType, typeSpecificData := findTypeOfData(data)
	fmt.Println("Data Type :\n", dataType)

	fmt.Println("Type specific data in bytes (INDEXING PURPOSE) :\n", typeSpecificData)
	//Response to client
	res := &document.DocumentTransferResponse{
		Response: "document recieved by server",
	}

	//variables that needs to be fed to database specific functions
	/*
		Data Field			Variable used		Type
		Database			database			string
		Collection			collection			string
		Namespace			namespace			string
		Data				data				map[string][]byte
		Indices				indices				[]string
		DataType			dataType			map[string]string
		Type Specific data	typeSpecificData	map[string][]byte  <- INDEXING PURPSE
	*/

	err:=store.InsertDocument(database,collection,namespace,data,indices)
	if err!=nil{
		return &document.DocumentTransferResponse{
			Response:             "",
		}, err
	}

	return res, nil
}

func main() {
	//create a new TiKV store
	err:=store.NewClient([]string{"127.0.0.1:2379"})
	if err!=nil{
		panic(err)
	}

	//read your env file and load them into ENV for this process
	err = godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	//change grpc server target from .env file
	grpcServerTarget := os.Getenv("GRPC_SERVER_TARGET")

	fmt.Println("-------------------Starting GRPC server-------------------")
	//just in case server crashes, get detailed log
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	//create a new listener
	lis, err := net.Listen("tcp", grpcServerTarget)

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
