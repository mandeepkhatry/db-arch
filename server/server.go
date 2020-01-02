package main

import (
	"context"
	"db-arch/pb/document"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type server struct{}

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

	res := &document.DocumentTransferResponse{
		Response: "document recieved by server",
	}
	//Call anish function to set to db

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

	s := grpc.NewServer()
	document.RegisterDocumentServiceServer(s, &server{})

	reflection.Register(s)

	err = s.Serve(lis)
	fmt.Println("check")

	if err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}

}
