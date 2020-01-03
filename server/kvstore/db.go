package kvstore

import (
	"errors"
	"log"
	"sync"

	"github.com/pingcap/tidb/config"
	"github.com/pingcap/tidb/store/tikv"
)

const (
	META_DBIDENTIFIER                 = "meta:dbidentifier"
	META_COLLECTIONIDENTIFIER         = "meta:collectionidentifier"
	META_NAMESPACEIDENTIFIER          = "meta:namespaceidentifier"
	META_DBID                         = "meta:dbid:"
	META_COLLECTIONID                 = "meta:collectionid:"
	META_NAMESPACEID                  = "meta:namespaceid:"
	META_DB                           = "meta:db:"
	META_COLLECTION                   = "meta:collection:"
	META_NAMESPACE                    = "meta:namespace:"
	DBIDENTIFIER_INITIALCOUNT         = uint16(1)
	COLLECTIONIDENTIFIER_INITIALCOUNT = uint32(1)
	NAMESPACEIDENTIFIER_INITIALCOUNT  = uint32(1)
)

//StoreClient is
type StoreClient struct {
	Client *tikv.RawKVClient
	m      sync.Mutex
}

//NewClient creates a new tikv.RawKVClient
func (s *StoreClient) NewClient(pdAddr []string) error {
	cli, err := tikv.NewRawKVClient([]string(pdAddr), config.Security{})
	if err != nil {
		return err
	}
	s.Client = cli
	return nil
}

//CloseClient closes tikv.RawKVClient
func (s *StoreClient) CloseClient() error {
	return s.Client.Close()
}

//Put inserts key,val to TiKV
func (s *StoreClient) Put(key []byte, value []byte) error {
	if len(key) == 0 {
		return errors.New("Key can't be empty")
	}
	err := s.Client.Put(key, value)
	if err != nil {
		return err
	}
	return nil
}

//PutBatch inserts key,val pairs in batch
//uses tikv.RawClient
func (s *StoreClient) PutBatch(args ...[]byte) error {
	keys := make([][]byte, 0)
	values := make([][]byte, 0)

	for i := 0; i < len(args); i += 2 {
		keys = append(keys, args[i])
		values = append(values, args[i+1])
	}

	log.Println("[[BatchPut]]")
	err := s.Client.BatchPut(keys, values)
	if err != nil {
		return err
	}
	return nil
}

//Get reads value for give key
func (s *StoreClient) Get(key []byte) ([]byte, error) {
	if len(key) == 0 {
		return nil, errors.New("key can't be empty")
	}
	val, err := s.Client.Get(key)
	if err != nil {
		return nil, err
	}
	return val, nil

}

//GetBatch retrieves values for given pair of keys in batch
func (s *StoreClient) GetBatch(keys [][]byte) ([][]byte, error) {
	if len(keys) == 0 {
		return nil, errors.New("keys are empty")
	}

	val, err := s.Client.BatchGet(keys)
	if err != nil {
		return nil, err
	}
	return val, nil
}

//DeleteKey deletes given key from TiKV
func (s *StoreClient) DeleteKey(key []byte) error {
	if len(key) == 0 {
		return errors.New("cannot delete empty key")
	}
	err := s.Client.Delete(key)
	if err != nil {
		return err
	}
	return nil
}

//DeleteKeyRange deletes key,val pairs from startKey to endKey from TiKV
func (s *StoreClient) DeleteKeyRange(startKey []byte, endKey []byte) error {
	if len(startKey) == 0 && len(endKey) == 0 {
		return errors.New("start or End keys cannot be empty")
	}
	err := s.Client.DeleteRange(startKey, endKey)
	if err != nil {
		return err
	}
	return nil
}

//Scan iterates from startKey to endKey upto within limit
func (s *StoreClient) Scan(startKey []byte, endKey []byte, limit int) ([][]byte, [][]byte, error) {
	keys := make([][]byte, 0)
	values := make([][]byte, 0)

	keys, values, err := s.Client.Scan(startKey, endKey, limit)
	if err != nil {
		return nil, nil, err
	}
	return keys, values, nil
}

//ReverseScan takes startKey, endKey and limit to scan in reverse direction in range [endKey,startKey)
//returns key,value,error
func (s *StoreClient) ReverseScan(startKey []byte, endKey []byte, limit int) ([][]byte, [][]byte, error) {
	keys := make([][]byte, 0)
	values := make([][]byte, 0)
	if len(startKey) == 0 {
		return nil, nil, errors.New("Can't scan from last without knowing startKey")
	}
	keys, values, err := s.Client.ReverseScan(startKey, endKey, limit)
	if err != nil {
		return nil, nil, err
	}
	return keys, values, nil
}
