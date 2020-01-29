package main

import (
	"context"
	"db-arch/model"
	"db-arch/pb/connection"
	"db-arch/pb/document"
	"db-arch/pb/query"
	"encoding/json"
	"net/http"
	"os"
	"strconv"
	"strings"

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
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(200)
			response := make(map[string]string)
			response["message"] = "document not created"
			json.NewEncoder(w).Encode(response)

			return
		}

		//Document object
		d := model.Document{
			Collection: dataInterface["collection"].(string),
			Data:       dataInterface["data"].(map[string]interface{}),
			Indices:    dataInterface["indices"].([]interface{}),
		}

		_, err = sendDocument(c, d)

		if err != nil {
			//TODO
			response := make(map[string]string)
			descFieldSplit := strings.Split(err.Error(), " desc = ")
			response["description"] = descFieldSplit[1]
			codeFieldSplit := strings.Split(descFieldSplit[0], " code = ")
			response["code"] = codeFieldSplit[1]

			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(200)
			json.NewEncoder(w).Encode(response)
			return

		}

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(201)
		response := make(map[string]string)
		response["message"] = "document created"
		json.NewEncoder(w).Encode(response)

		return

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

		d := model.Query{
			Query: dataInterface["query"].(string),
		}

		res, err := sendQuery(c, d)

		if err != nil {
			response := make(map[string]string)
			descFieldSplit := strings.Split(err.Error(), " desc = ")
			response["description"] = descFieldSplit[1]
			codeFieldSplit := strings.Split(descFieldSplit[0], " code = ")
			response["code"] = codeFieldSplit[1]

			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(200)
			json.NewEncoder(w).Encode(response)
			return

		}

		if len(res.GetResponse()) > 0 {
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(200)

			result := make(map[string](map[string]interface{}))

			//resultInterface represents different types of data
			var resultInterface interface{}

			for responseKey, responseValue := range res.GetResponse() {
				eachKV := make(map[string]interface{})
				for fieldName, fieldValue := range responseValue.GetResult() {
					print("fieldname, fieldvalue : ", fieldName, fieldValue)
					json.Unmarshal(fieldValue, &resultInterface)
					eachKV[fieldName] = resultInterface
				}
				//key as string and value as map[string]interface{}
				result["result "+strconv.Itoa(responseKey)] = eachKV
			}

			json.NewEncoder(w).Encode(result)

			return

		}
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(200)
		response := make(map[string]string)
		response["message"] = "query not found"
		json.NewEncoder(w).Encode(response)
		return
	}
}

func connectDatabase(c connection.ConnectionServiceClient) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		dataInterface := make(map[string]interface{})

		err := json.NewDecoder(r.Body).Decode(&dataInterface)

		if err != nil {
			panic(err)
		}

		d := model.Connection{
			Database:  dataInterface["database"].(string),
			Namespace: dataInterface["namespace"].(string),
		}
		res, err := sendConnection(c, d)

		response := make(map[string]string)
		if err != nil {
			response["connection_status"] = "failed to establish connection"
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(200)
			json.NewEncoder(w).Encode(response)
			return
		}

		response["connection_status"] = res.GetResponse()
		w.Header().Add("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return

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

	c3 := connection.NewConnectionServiceClient(conn)

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
	router.HandleFunc("/connection", connectDatabase(c3)).Methods("POST")
	http.ListenAndServe(":8000", router)

}

//unary API call to send document to server.go
func sendDocument(c document.DocumentServiceClient, d model.Document) (*document.DocumentTransferResponse, error) {
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
			Collection: d.Collection,
			Data:       newData,
			Indices:    indices,
		},
	}

	res, err := c.DocumentTransfer(context.Background(), req)
	if err != nil {
		return res, err
	}

	fmt.Println("----------------RESPONSE----------------")
	return res, nil
}

//unary API call to send query to server.go
func sendQuery(c query.QueryServiceClient, d model.Query) (*query.QueryTransferResponse, error) {
	fmt.Println("----------------Sending Query-------------------")

	//Document Transfer Request
	req := &query.QueryTransferRequest{
		Request: &query.Query{
			Query: d.Query,
		},
	}
	res, err := c.QueryTransfer(context.Background(), req)
	if err != nil {
		return res, err
	}

	fmt.Println("----------------QUERY RESPONSE----------------")
	fmt.Println("Query recieved : ", res)
	return res, nil
}

//unary API call to send connection to server.go
func sendConnection(c connection.ConnectionServiceClient, d model.Connection) (*connection.ConnectionTransferResponse, error) {
	fmt.Println("-----------------Establishing Connection-------------")
	fmt.Println("DB N : ", d.Database, d.Namespace)
	//Connection Transfer Request
	req := &connection.ConnectionTransferRequest{
		Request: &connection.Connection{
			Database:  d.Database,
			Namespace: d.Namespace,
		},
	}

	res, err := c.ConnectionTransfer(context.Background(), req)
	if err != nil {
		return res, err
	}

	fmt.Println("----------CONNECTION RESPONSE------------")
	fmt.Println("Connection recieved : ", res)
	return res, nil

}
