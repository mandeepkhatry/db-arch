package main

import (
	"context"
	"db-arch/model"
	"db-arch/pb/document"
	"encoding/json"
	"net/http"

	"fmt"
	"log"

	"github.com/gorilla/mux"
	"google.golang.org/grpc"
)

//Post Handler Function
func postDocument(logger chan model.Document) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("inside")

		data := make(map[string]interface{})

		err := json.NewDecoder(r.Body).Decode(&data)

		if err != nil {
			panic(err)
		}
		d := model.Document{
			Database:   data["database"].(string),
			Collection: data["collection"].(string),
			Namespace:  data["namespace"].(string),
			Data:       data["data"],
		}
		fmt.Println(d)
		logger <- d
	}
}

func main() {

	//handlerChannel represents a channel from which documents are passed via HTTP request
	handlerChannel := make(chan model.Document, 10)

	noExit := make(chan string)

	//Starting GRPC Client
	fmt.Println("-------------------Starting GRPC client-------------------")
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())

	//Register client
	c := document.NewDocumentServiceClient(conn)

	if err != nil {
		log.Fatalf("Client couldn't connect to server %v", err)
		panic(err)
	}

	//close connection when program closes
	defer conn.Close()

	//Go routine to recieve documents from HTTP Request
	go func() {
		for {
			val, ok := <-handlerChannel
			if ok {
				sendDocument(c, val)
			}
		}

	}()

	//Starting server
	router := mux.NewRouter()

	fmt.Println("-------------------Starting client-------------------")
	router.HandleFunc("/documents", postDocument(handlerChannel)).Methods("POST")
	http.ListenAndServe(":8000", router)

	<-noExit

}

//unary API call to send document to server.go
func sendDocument(c document.DocumentServiceClient, d model.Document) {
	fmt.Println("----------------Sending Document-------------------")

	data, err := json.Marshal(d.Data)

	if err != nil {
		log.Fatal(err)
		panic(err)
	}

	//Document Transfer Request
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

	fmt.Println("----------------RESPONSE----------------")
	fmt.Println(res)
}
