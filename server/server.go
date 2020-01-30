package main

import (
	"context"
	"db-arch/pb/connection"
	"db-arch/pb/document"
	"db-arch/pb/query"
	"db-arch/server/internal/def"
	"db-arch/server/internal/engine"
	"db-arch/server/internal/engine/parser"
	"db-arch/server/internal/kvstore"
	"db-arch/server/io"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

//create Store interface
var store io.Store

//Engine
var eng engine.Engine

//Server struct
type server struct{}

//DocumentTransfer function that recieves request from client and returns response
func (*server) DocumentTransfer(ctx context.Context, req *document.DocumentTransferRequest) (*document.DocumentTransferResponse, error) {
	fmt.Println("-----Got Document as Request-----")

	if eng.DBName == "" {
		//Response to client
		res := &document.DocumentTransferResponse{
			Response: "connect to database",
		}
		return res, def.CONNECTION_NOT_ESTABLISHED
	}

	collection := req.GetRequest().GetCollection()

	//data in form of bytes
	data := req.GetRequest().GetData()
	indices := req.GetRequest().GetIndices()

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

	err := eng.InsertDocument(store, collection, data, indices)
	if err != nil {
		statusCode := def.ERRTYPE[err]
		return &document.DocumentTransferResponse{
			Response: "",
		}, status.Error(statusCode, err.Error())
	}

	return res, nil
}

//DocumentTransfer function that recieves request from client and returns response
func (*server) QueryTransfer(ctx context.Context, req *query.QueryTransferRequest) (*query.QueryTransferResponse, error) {
	fmt.Println("-----Got Query as Request-----")

	if eng.DBName == "" {
		//Response to client
		res := &query.QueryTransferResponse{}
		return res, def.CONNECTION_NOT_ESTABLISHED
	}

	rawQuery := req.GetRequest().GetQuery()
	print("Recieved raw query :", rawQuery)

	collection, postfixQuery, err := parser.ParseQuery(rawQuery)
	fmt.Println("[[server.go/collection]]", collection)
	fmt.Println("[[server.go/postfixQuery]]", postfixQuery)
	if err != nil {
		res := &query.QueryTransferResponse{}
		return res, err
	}

	//TODO SearchDocumet contains code of evalation from postfixQuery
	resultArray, err := eng.SearchDocument(store, collection, postfixQuery)

	if err != nil {
		statusCode := def.ERRTYPE[err]
		return &query.QueryTransferResponse{}, status.Error(statusCode, err.Error())
	}
	fmt.Println("[[server.go]]resultArray", resultArray)
	var resultInBytes map[string][]byte
	response := make([]*query.Response, 0)

	for _, v := range resultArray {
		//bytes to map[string][]byte
		json.Unmarshal(v, &resultInBytes)
		var each_response query.Response

		each_response.Result = resultInBytes
		response = append(response, &each_response)
	}

	//Response to client
	res := &query.QueryTransferResponse{
		Response: response,
	}

	return res, nil
}

func (*server) ConnectionTransfer(ctx context.Context, req *connection.ConnectionTransferRequest) (*connection.ConnectionTransferResponse, error) {
	fmt.Println("--------Establishing connection--------")
	database := req.GetRequest().GetDatabase()
	namespace := req.GetRequest().GetNamespace()
	print("DATABASE , NAMESPACE : ", database, namespace)

	eng.DBName = database
	eng.Namespace = namespace
	print("----------ConnectDB function calling-----------")
	err := eng.ConnectDB(store)
	if err != nil {
		return &connection.ConnectionTransferResponse{}, status.Error(codes.Aborted, err.Error())
	}

	res := &connection.ConnectionTransferResponse{
		Response: "connection established with " + database + ":" + namespace,
	}
	return res, nil
}
func main() {
	//create a new badger store from factory
	store = kvstore.NewBadgerFactory([]string{}, "./data/badger")

	//create tikv
	// store=kvstore.NewTiKVFactory([]string{"addr here"},"")

	//read your env file and load them into ENV for this process
	err := godotenv.Load()
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
	connection.RegisterConnectionServiceServer(s, &server{})

	reflection.Register(s)

	err = s.Serve(lis)

	if err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}

}
