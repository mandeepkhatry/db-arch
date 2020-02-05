# db-arch

## Command to generate .pb.go file
```
protoc pb/document/document.proto --go_out=plugins=grpc:.
protoc pb/query/query.proto --go_out=plugins=grpc:.

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


**Build project**
```
make build 
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
```

**Command to run tikv pd**
```
./bin/pd-server --data-dir=pd --log-file=pd.log
```

**Command to run tikv server**
```
./bin/tikv-server --pd="127.0.0.1:2379" --data-dir=tikv --log-file=tikv.log
```


*Note: use root permission to run tikv cluster*
