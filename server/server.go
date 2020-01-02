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

//TODO buildData function for type assertion
//Status : Incomplete

func buildData(data []byte) {

	var data_interface map[string]interface{}

	err := json.Unmarshal(data, &data_interface)

	if err != nil {
		panic(err)
	}

	for k, v := range data_interface {
		fmt.Println(k, v)

	}

	print(data_interface)

}

//Server struct
type server struct{}

//DocumentTransfer function that recieves request from client and returns response
func (*server) DocumentTransfer(ctx context.Context, req *document.DocumentTransferRequest) (*document.DocumentTransferResponse, error) {
	fmt.Println("-----Got Document as Request-----")

	database := req.GetRequest().GetDatabase()
	collection := req.GetRequest().GetCollection()
	namespace := req.GetRequest().GetNamespace()
	data := req.GetRequest().GetData()

	fmt.Println("Document : ", database)
	fmt.Println("Collection : ", collection)
	fmt.Println("Namespace : ", namespace)
	fmt.Println("Data : ", data)

	//Response to client
	res := &document.DocumentTransferResponse{
		Response: "document recieved by server",
	}

	//Call function to set to db with data we have

	return res, nil
}

func main() {
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
	fmt.Println("check")

	if err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}

}
