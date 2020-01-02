package main

import (
	"context"
	"db-arch/model"
	"db-arch/pb/document"
	"encoding/json"

	"fmt"
	"google.golang.org/grpc"
	"log"
)

func main() {

	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())

	if err != nil {
		log.Fatalf("Client couldn't connect to server %v", err)
		panic(err)
	}
	//close connection when program closes
	defer conn.Close()
	//register client
	c := document.NewDocumentServiceClient(conn)

	data := make(map[string]string)

	data["Name"] = "Mandeep"

	d := model.Document{
		Database:   "db1",
		Collection: "c1",
		Namespace:  "n1",
		Data:       data,
	}

	sendDocument(c, d)

}

//unary API
func sendDocument(c document.DocumentServiceClient, d model.Document) {
	fmt.Println("----------------Send Document-------------------")

	data, err := json.Marshal(d.Data)

	if err != nil {
		log.Fatal(err)
		panic(err)
	}

	req := &document.DocumentTransferRequest{
		Request: &document.Document{
			Database:   d.Database,
			Collection: d.Collection,
			Namespace:  d.Namespace,
			Data:       data,
		},
	}

	res, err := c.DocumentTransfer(context.Background(), req)
	if err != nil {
		log.Fatalf("%v", err)
	}

	fmt.Println("RESPONSE----- %v", res)
}
