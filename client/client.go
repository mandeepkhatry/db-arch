package main

import (
	"bytes"
	"context"
	"db-arch/model"
	"db-arch/pb/connection"
	"db-arch/pb/document"
	"db-arch/pb/query"
	"encoding/gob"
	"os"
	"strings"
	"sync"

	"fmt"
	"log"

	"github.com/joho/godotenv"
	"github.com/nats-io/nats.go"
	"google.golang.org/grpc"
)

// //Post Handler Function
// func postDocument(c document.DocumentServiceClient) func(http.ResponseWriter, *http.Request) {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		dataInterface := make(map[string]interface{})

// 		decoderInstance := json.NewDecoder(r.Body)
// 		decoderInstance.UseNumber()
// 		err := decoderInstance.Decode(&dataInterface)

// 		var postResponse response.PostResponse
// 		var metaResponse response.PostMetaResponse

// 		if err != nil {
// 			w.Header().Add("Content-Type", "application/json")
// 			w.WriteHeader(404)
// 			postResponse.Message = "document not created"
// 			metaResponse.Status = false
// 			metaResponse.Description = "not created"
// 			metaResponse.Code = ""
// 			postResponse.Metadata = metaResponse
// 			json.NewEncoder(w).Encode(postResponse)

// 			return
// 		}

// 		//Document object
// 		d := model.Document{
// 			Collection: dataInterface["collection"].(string),
// 			Data:       dataInterface["data"].(map[string]interface{}),
// 			Indices:    dataInterface["indices"].([]interface{}),
// 		}

// 		_, err = sendDocument(c, d)

// 		if err != nil {
// 			descFieldSplit := strings.Split(err.Error(), " desc = ")
// 			codeFieldSplit := strings.Split(descFieldSplit[0], " code = ")

// 			metaResponse.Description = descFieldSplit[1]
// 			metaResponse.Code = codeFieldSplit[1]
// 			metaResponse.Status = false

// 			postResponse.Message = "document not created"
// 			postResponse.Metadata = metaResponse

// 			w.Header().Add("Content-Type", "application/json")
// 			w.WriteHeader(404)
// 			json.NewEncoder(w).Encode(postResponse)
// 			return

// 		}

// 		w.Header().Add("Content-Type", "application/json")
// 		w.WriteHeader(201)

// 		postResponse.Message = "document created"
// 		metaResponse.Code = ""
// 		metaResponse.Status = true
// 		metaResponse.Description = "created"
// 		postResponse.Metadata = metaResponse

// 		json.NewEncoder(w).Encode(postResponse)

// 		return

// 	}
// }

// //Query Handler Function
// func queryDocument(c query.QueryServiceClient) func(http.ResponseWriter, *http.Request) {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		dataInterface := make(map[string]interface{})

// 		err := json.NewDecoder(r.Body).Decode(&dataInterface)

// 		if err != nil {
// 			panic(err)
// 		}

// 		d := model.Query{
// 			Query: dataInterface["query"].(string),
// 		}

// 		res, err := sendQuery(c, d)

// 		var queryResponse response.QueryResponse

// 		if err != nil {
// 			queryResponse.Result = nil
// 			descFieldSplit := strings.Split(err.Error(), " desc = ")
// 			queryResponse.Metadata.Description = descFieldSplit[1]
// 			codeFieldSplit := strings.Split(descFieldSplit[0], " code = ")
// 			queryResponse.Metadata.Code = codeFieldSplit[1]
// 			queryResponse.Metadata.Status = false
// 			w.Header().Add("Content-Type", "application/json")
// 			w.WriteHeader(404)
// 			json.NewEncoder(w).Encode(queryResponse)
// 			return

// 		}

// 		if len(res.GetResponse()) > 0 {
// 			w.Header().Add("Content-Type", "application/json")
// 			w.WriteHeader(200)

// 			//resultInterface represents different types of data
// 			var resultInterface interface{}
// 			var results = make([]map[string]interface{}, 0)

// 			for _, responseValue := range res.GetResponse() {
// 				var tempResult = make(map[string]interface{})
// 				for fieldName, fieldValue := range responseValue.GetResult() {
// 					json.Unmarshal(fieldValue, &resultInterface)
// 					tempResult[fieldName] = resultInterface
// 				}
// 				results = append(results, tempResult)

// 				//key as string and value as map[string]interface{}
// 			}
// 			//meta data
// 			queryResponse.Metadata.Status = true
// 			queryResponse.Metadata.Code = ""
// 			queryResponse.Metadata.Description = ""

// 			queryResponse.Result = results

// 			json.NewEncoder(w).Encode(queryResponse)

// 			return

// 		}
// 		w.Header().Add("Content-Type", "application/json")
// 		w.WriteHeader(404)
// 		queryResponse.Result = nil
// 		queryResponse.Metadata.Description = "results not found"
// 		queryResponse.Metadata.Code = ""
// 		queryResponse.Metadata.Status = false
// 		json.NewEncoder(w).Encode(queryResponse)
// 		return
// 	}
// }

// func connectDatabase(c connection.ConnectionServiceClient) func(http.ResponseWriter, *http.Request) {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		dataInterface := make(map[string]interface{})

// 		err := json.NewDecoder(r.Body).Decode(&dataInterface)

// 		if err != nil {
// 			panic(err)
// 		}

// 		d := model.Connection{
// 			Database:  dataInterface["database"].(string),
// 			Namespace: dataInterface["namespace"].(string),
// 		}
// 		res, err := sendConnection(c, d)

// 		var connectionResponse response.ConnectionResponse

// 		if err != nil {
// 			connectionResponse.Message = "failed to establish connection"
// 			connectionResponse.Metadata.Status = false
// 			w.Header().Add("Content-Type", "application/json")
// 			w.WriteHeader(200)
// 			json.NewEncoder(w).Encode(connectionResponse)
// 			return
// 		}

// 		connectionResponse.Message = res.GetResponse()
// 		connectionResponse.Metadata.Status = true
// 		w.Header().Add("Content-Type", "application/json")
// 		json.NewEncoder(w).Encode(connectionResponse)
// 		return

// 	}
// }

var nc *nats.Conn

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
	//set 20MB max message size for grpc
	conn, err := grpc.Dial(grpcServerTarget, grpc.WithInsecure(), grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(20*1024*1024), grpc.MaxCallSendMsgSize(512*1024*1024)))

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

	fmt.Println("Connected to NATS cluster")
	nc, err = nats.Connect("0.0.0.0:4222")

	nc.Subscribe("db.*", func(m *nats.Msg) {
		subject := m.Subject
		nextSub := strings.Split(subject, ".")[1]
		fmt.Println("INSIDE DB")
		gob.Register(map[string]interface{}{})
		gob.Register([]interface{}{})
		gob.Register([]map[string]interface{}{})

		buf := bytes.NewBuffer(m.Data)
		dec := gob.NewDecoder(buf)
		requestData := make(map[string]interface{})
		err := dec.Decode(&requestData)
		if err != nil {
			nc.Publish(m.Reply, []byte(strings.Split(err.Error(), "desc = ")[1]))
		}

		switch nextSub {
		case "insertdocument":
			indices := make([]string, 0)
			for _, v := range requestData["payload"].(map[string]interface{})["indices"].([]interface{}) {
				indices = append(indices, v.(string))
			}

			data := make(map[string][]byte, 0)
			for key, value := range requestData["payload"].(map[string]interface{})["data"].(map[string]interface{}) {
				var buf bytes.Buffer
				enc := gob.NewEncoder(&buf)
				if err := enc.Encode(value); err != nil {
					nc.Publish(m.Reply, []byte(strings.Split(err.Error(), "desc = ")[1]))
				}
				data[key] = buf.Bytes()
			}

			doc := document.Document{
				Collection: (requestData["payload"].(map[string]interface{})["collection"]).(string),
				Data:       data,
				Indices:    indices,
			}
			//doc := document.Document{}
			//err := proto.Unmarshal(m.Data, &doc)
			//if err != nil {
			//	panic(err)
			//}
			fmt.Println("DOC : ", doc)
			res, err := sendDocument(c1, doc)

			fmt.Println("RESPONSE : ", res)

			if err != nil {
				desc := strings.Split(err.Error(), " desc = ")[1]
				nc.Publish(m.Reply, []byte(desc))
			} else {
				var buf bytes.Buffer
				enc := gob.NewEncoder(&buf)
				data := make(map[string]interface{})
				data["response"] = res.GetResponse()
				if err := enc.Encode(data); err != nil {
					log.Fatal(err)
				}
				nc.Publish(m.Reply, buf.Bytes())
			}
		case "connection":

			con := connection.Connection{
				Database:  (requestData["payload"].(map[string]interface{})["database"]).(string),
				Namespace: (requestData["payload"].(map[string]interface{})["namespace"]).(string),
			}

			//err := proto.Unmarshal(m.Data, &con)
			//if err != nil {
			//	panic(err)
			//}

			data := model.Connection{
				Database:  con.GetDatabase(),
				Namespace: con.GetNamespace(),
			}
			fmt.Println("received connection request")
			res, err := sendConnection(c3, data)

			if err != nil {
				nc.Publish(m.Reply, []byte("failed to establish connection"))
			} else {
				var buf bytes.Buffer
				enc := gob.NewEncoder(&buf)
				data := make(map[string]interface{})
				data["response"] = res.GetResponse()
				if err := enc.Encode(data); err != nil {
					log.Fatal(err)
				}
				nc.Publish(m.Reply, buf.Bytes())
			}

		case "querydocument":
			inquery := query.Query{
				Query: (requestData["payload"].(map[string]interface{})["query"]).(string),
			}
			//err := proto.Unmarshal(m.Data, &inquery)
			//if err != nil {
			//	panic(err)
			//}
			//rawquery := query.Query{
			//	Query: inquery.GetQuery(),
			//}

			res, err := sendQuery(c2, inquery)

			if err != nil {
				nc.Publish(m.Reply, []byte(strings.Split(err.Error(), "desc = ")[1]))
			} else {
				response := make(map[string]interface{})
				var results = make([]map[string]interface{}, 0)
				var resultInterface interface{}

				if len(res.GetResponse()) > 0 {
					fmt.Println("QUERY RESPONSE : ", res.GetResponse())

					for _, responseValue := range res.GetResponse() {
						var tempResult = make(map[string]interface{})
						for fieldName, fieldValue := range responseValue.GetResult() {
							buf := bytes.NewBuffer(fieldValue)
							dec := gob.NewDecoder(buf)
							dec.Decode(&resultInterface)
							tempResult[fieldName] = resultInterface
						}
						results = append(results, tempResult)

					}
					response["data"] = results
					response["status"] = "found"

				} else {
					response["data"] = results
					response["status"] = "no query"
				}
				fmt.Println("RESPONSE : ")
				var buf bytes.Buffer
				enc := gob.NewEncoder(&buf)

				if err := enc.Encode(response); err != nil {
					log.Fatal(err)
				}
				nc.Publish(m.Reply, buf.Bytes())
			}

		}

	})
	wg := sync.WaitGroup{}
	wg.Add(1)
	wg.Wait()
	//Starting server
	// router := mux.NewRouter()

	// fmt.Println("-------------------Starting client-------------------")
	// router.HandleFunc("/documents", postDocument(c1)).Methods("POST")
	// router.HandleFunc("/query", queryDocument(c2)).Methods("POST")
	// router.HandleFunc("/connection", connectDatabase(c3)).Methods("POST")
	// http.ListenAndServe(":8000", router)

}

//unary API call to send document to server.go
func sendDocument(c document.DocumentServiceClient, d document.Document) (*document.DocumentTransferResponse, error) {
	fmt.Println("----------------Sending Document-------------------")

	//Document Transfer Request
	req := &document.DocumentTransferRequest{
		Request: &document.Document{
			Collection: d.GetCollection(),
			Data:       d.GetData(),
			Indices:    d.GetIndices(),
		},
	}

	res, err := c.DocumentTransfer(context.Background(), req)
	fmt.Println("ERROR : ", err)
	if err != nil {
		return res, err
	}

	fmt.Println("----------------RESPONSE----------------")
	fmt.Println("RESPONSE : ", res)
	return res, nil
}

//unary API call to send query to server.go
func sendQuery(c query.QueryServiceClient, d query.Query) (*query.QueryTransferResponse, error) {
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
	//fmt.Println("Query recieved : ", res)
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
