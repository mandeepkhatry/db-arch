package main

import (
	"context"
	"db-arch/pb/document"
	"db-arch/pb/query"
	"fmt"
	"log"
	"net"
	"os"

	"db-arch/server/kvstore"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var store kvstore.StoreClient

//TODO Business Logic to map golang specific type

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
		Type Specific data	typeSpecificData	map[string][]byte  <- INDEXING PURPOSE
	*/

	err := store.InsertDocument(database, collection, namespace, data, indices)
	if err != nil {
		return &document.DocumentTransferResponse{
			Response: "",
		}, err
	}

	return res, nil
}

//DocumentTransfer function that recieves request from client and returns response
func (*server) QueryTransfer(ctx context.Context, req *query.QueryTransferRequest) (*query.QueryTransferResponse, error) {
	fmt.Println("-----Got Query as Request-----")

	database := req.GetRequest().GetDatabase()
	collection := req.GetRequest().GetCollection()
	namespace := req.GetRequest().GetNamespace()
	//data in form of bytes
	queryData := req.GetRequest().GetQuerydata()
	/*
		To unmarshal data use this :
			var dataInterface map[string]interface{}
			err := json.Unmarshal(data, &dataInterface)

			if err != nil {
				panic(err)
			}
	*/

	fmt.Println("Database : ", database)
	fmt.Println("Collection : ", collection)
	fmt.Println("Namespace : ", namespace)
	fmt.Println("Query Data : ", queryData)

	_, err := store.SearchDocument(database, collection, namespace, queryData)
	if err != nil {
		return &query.QueryTransferResponse{}, err
	}
	
	//Response to client
	res := &query.QueryTransferResponse{
		Response: "query recieved by server",
	}

	return res, nil
}

func main() {
	//create a new TiKV store
	err := store.NewClient([]string{"127.0.0.1:2379"})
	if err != nil {
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
	query.RegisterQueryServiceServer(s, &server{})

	reflection.Register(s)

	err = s.Serve(lis)

	if err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}

}
