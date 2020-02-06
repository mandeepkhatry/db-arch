# db-arch

## Command to generate .pb.go file
---
```
protoc pb/document/document.proto --go_out=plugins=grpc:.
protoc pb/query/query.proto --go_out=plugins=grpc:.
protoc pb/connection/connection.proto --go_out=plugins=grpc:.

```

## Run two terminals for server and client
---

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
	"collection" : "sectionB",
	"data" :{
		"name" : "Mandeep Khatry",
		"age" : 26,
		"isEmployee" : true,
		"salary"	: 55,
		"joined_date":"2019-08-08",
		"time" : "2020-02-07",
		"address":"Dhobighat",
		"school":"institute of engineering",
		"company":"new rara labs"
	},
	"indices" : ["name", "age", "salary" ,"isEmployee,joined_date"]
}

Single indices : name, age, salary
Compound indices : (isEmployee, joined_date)
```
-----
```/query```  
``` 
Methods: POST

Request

{
		"query" : "@sectionB (name=\"Anish Bhusal\" OR age>=26 OR address=\"ktm\") AND (joined_date>\"2010-09-23\" OR company=\"gokyo labs\"")
}

Supported Operators : AND, OR, >, <, <=, >=, !=, =

Note : String Format (enclosed by \" and \")
```

-----
```/connection```  
``` 
Methods: POST

Request

{
	"database":"school",
	"namespace":"ideal-model"
}

```

## KVStores implemented
---
```
BadgerDB
TiKV
```


## Tikv

----

**Command to run tikv pd**
```
./bin/pd-server --data-dir=pd --log-file=pd.log
```

**Command to run tikv server**
```
./bin/tikv-server --pd="127.0.0.1:2379" --data-dir=tikv --log-file=tikv.log
```


*Note: use root permission to run tikv cluster*
