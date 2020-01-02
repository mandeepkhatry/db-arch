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
		fmt.Println("here")
	}
}

func main() {

	handlerChannel := make(chan model.Document, 10)

	no_exit := make(chan string)

	fmt.Println("Checking")

	fmt.Println("Starting GRPC client ...")
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())

	//register client
	c := document.NewDocumentServiceClient(conn)

	if err != nil {
		log.Fatalf("Client couldn't connect to server %v", err)
		panic(err)
	}
	//close connection when program closes
	defer conn.Close()

	go func() {
		fmt.Println("Here")
		for {
			val, ok := <-handlerChannel
			if ok {
				fmt.Println("Here from channel : ", val)

				sendDocument(c, val)
			}
		}

	}()

	router := mux.NewRouter()

	fmt.Println("Starting client ...")

	router.HandleFunc("/documents", postDocument(handlerChannel)).Methods("POST")

	http.ListenAndServe(":8000", router)

	<-no_exit

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
