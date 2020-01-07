package main

import (
	"context"
	"db-arch/model"
	"db-arch/pb/document"
	"db-arch/pb/query"
	"encoding/json"
	"net/http"
	"os"

	"fmt"
	"log"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
)

//Post Handler Function
func postDocument(c document.DocumentServiceClient) func(http.ResponseWriter, *http.Request) {
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
			Data:       dataInterface["data"].(map[string]interface{}),
			Indices:    dataInterface["indices"].([]interface{}),
		}
		fmt.Println("Document recieved : ", d)
		sendDocument(c, d)
	}
}

//Query Handler Function
func queryDocument(c query.QueryServiceClient) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		dataInterface := make(map[string]interface{})

		err := json.NewDecoder(r.Body).Decode(&dataInterface)

		if err != nil {
			panic(err)
		}

		//Document object
		d := model.Query{
			Database:   dataInterface["database"].(string),
			Collection: dataInterface["collection"].(string),
			Namespace:  dataInterface["namespace"].(string),
			Querydata:  dataInterface["data"].(map[string]interface{}),
		}
		fmt.Println("Document recieved : ", d)
		sendQuery(c, d)
	}
}
func main() {
	//read your env file and load them into ENV for this process
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	//change grpc server target from .env file
	grpcServerTarget := os.Getenv("GRPC_SERVER_TARGET")

	//Starting GRPC Client
	fmt.Println("-------------------Starting GRPC client-------------------")
	conn, err := grpc.Dial(grpcServerTarget, grpc.WithInsecure())

	//Register client
	c1 := document.NewDocumentServiceClient(conn)

	c2 := query.NewQueryServiceClient(conn)

	if err != nil {
		log.Fatalf("Client couldn't connect to server %v", err)
		panic(err)
	}

	//close connection when program closes
	defer conn.Close()

	//Starting server
	router := mux.NewRouter()

	fmt.Println("-------------------Starting client-------------------")
	router.HandleFunc("/documents", postDocument(c1)).Methods("POST")
	router.HandleFunc("/query", queryDocument(c2)).Methods("POST")
	http.ListenAndServe(":8000", router)

}

//unary API call to send document to server.go
func sendDocument(c document.DocumentServiceClient, d model.Document) {
	fmt.Println("----------------Sending Document-------------------")
	indices := make([]string, 0)

	for _, v := range d.Indices {
		indices = append(indices, v.(string))
	}

	newData := make(map[string][]byte)

	for k, v := range d.Data {

		bytedata, err := json.Marshal(v)
		if err != nil {
			panic(err)
		}
		newData[k] = bytedata
	}

	//Document Transfer Request
	req := &document.DocumentTransferRequest{
		Request: &document.Document{
			Database:   d.Database,
			Collection: d.Collection,
			Namespace:  d.Namespace,
			Data:       newData,
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

//unary API call to send query to server.go
func sendQuery(c query.QueryServiceClient, d model.Query) {
	fmt.Println("----------------Sending Query-------------------")

	newData := make(map[string][]byte)

	for k, v := range d.Querydata {

		bytedata, err := json.Marshal(v)
		if err != nil {
			panic(err)
		}
		newData[k] = bytedata
	}

	//Document Transfer Request
	req := &query.QueryTransferRequest{
		Request: &query.Query{
			Database:   d.Database,
			Collection: d.Collection,
			Namespace:  d.Namespace,
			Querydata:  newData,
		},
	}
	res, err := c.QueryTransfer(context.Background(), req)
	if err != nil {
		log.Fatalf("%v", err)
	}

	fmt.Println("----------------QUERY RESPONSE----------------")
	fmt.Println(res)
}
