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
		dataInterface := make(map[string]interface{})

		err := json.NewDecoder(r.Body).Decode(&dataInterface)

		if err != nil {
			panic(err)
		}

		//Document object
		d := model.Document{
			Database:   dataInterface["database"].(string),
			Collection: dataInterface["collection"].(string),
			Namespace:  dataInterface["namespace"].(string),
			Data:       dataInterface["data"],
			Indices:    dataInterface["indices"].([]interface{}),
		}
		fmt.Println("Document recieved : ", d)
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

	close(handlerChannel)
	close(noExit)

}

//unary API call to send document to server.go
func sendDocument(c document.DocumentServiceClient, d model.Document) {
	fmt.Println("----------------Sending Document-------------------")

	data, err := json.Marshal(d.Data)

	if err != nil {
		log.Fatal(err)
		panic(err)
	}

	indices := make([]string, 0)

	for _, v := range d.Indices {
		indices = append(indices, v.(string))
	}

	//Document Transfer Request
	req := &document.DocumentTransferRequest{
		Request: &document.Document{
			Database:   d.Database,
			Collection: d.Collection,
			Namespace:  d.Namespace,
			Data:       data,
			Indices:    indices,
		},
	}

	res, err := c.DocumentTransfer(context.Background(), req)
	if err != nil {
		log.Fatalf("%v", err)
	}

	fmt.Println("----------------RESPONSE----------------")
	fmt.Println(res)
}
