# db-arch

## Command to generate .pb.go file
```
protoc pb/document/document.proto --go_out=plugins=grpc:.

```

## Run two terminals for server and client

### Server (command)
```
go run server/server.go

```

### Client (command)
```
go run client/client.go

```

## API Endpoints
-----
```/documents```  
``` 
Methods: POST

Request

{
	"database" : "db1",
	"collection" : "c1",
	"namespace" : "n1",
	"data" :{
		"name" : "Mandeep",
		"age" : 19.0,
		"salary" :19000.5,
		"isEmployee" : true
	},
	"indices" : ["name","age"]
}
