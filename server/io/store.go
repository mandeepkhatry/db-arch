package io

//Store interface
type Store interface{
	//NewClient creates a new db client
	//takes placement driver address as arg
	NewClient(pdAddr string) (error)
	//CloseClient closes DB client
	CloseClient()error
	//Put inserts key,val to DB
	Put(key []byte,value []byte) error
	//PutBatch inserts key,val pairs in batch
	PutBatch(keys [][]byte,values [][]byte) error
	//Get retrieves value for given key
	Get(key []byte)([]byte, error)
	//GetBatch retrieves values for given collection of keys in batch
	GetBatch(keys [][]byte)([][]byte,error)
	//DeleteKey deletes given key from DB
	DeleteKey(key []byte)error
	//DeleteKeyRange deletes key,val pairs from startKey to endKey from DB
	DeleteKeyRange(startKey []byte,endKey []byte)error
	//Scan iterates from startKey to endKey for given limit
	//leaving limit empty throws error
	Scan(startKey []byte,endKey []byte,limit int)([][]byte, [][]byte, error)
	//ReverseScan takes startKey,endKey and limit to scan in reverse direction
	//startKey can't be set to "" empty
	ReverseScan(startKey []byte, endKey []byte, limit int) ([][]byte, [][]byte, error)
}
