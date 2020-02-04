package tikvdb

import (
	"db-arch/server/internal/def"
	"log"
	"sync"

	"github.com/pingcap/tidb/config"
	"github.com/pingcap/tidb/store/tikv"
)

//StoreClient implements tikvdb.RawKVClient
type StoreClient struct {
	Client *tikv.RawKVClient
	m      sync.Mutex
}

//NewClient creates a new tikvdb.RawKVClient
func (s *StoreClient) NewClient(pdAddr []string, dbDIR string) error {
	cli, err := tikv.NewRawKVClient([]string(pdAddr), config.Security{})
	if err != nil {
		return err
	}
	s.Client = cli
	return nil
}

//CloseClient closes tikvdb.RawKVClient
func (s *StoreClient) CloseClient() error {
	return s.Client.Close()
}

//Put inserts key,val to TiKV
func (s *StoreClient) Put(key []byte, value []byte) error {
	if len(key) == 0 {
		return def.KEY_EMPTY
	}
	err := s.Client.Put(key, value)
	if err != nil {
		return err
	}
	return nil
}

//PutBatch inserts key,val pairs in batch
//uses tikvdb.RawClient
func (s *StoreClient) PutBatch(keys [][]byte, values [][]byte) error {
	log.Println("[[BatchPut]]")
	err := s.Client.BatchPut(keys, values)
	if err != nil {
		return err
	}
	return nil
}

//Get reads value for given key
func (s *StoreClient) Get(key []byte) ([]byte, error) {
	if len(key) == 0 {
		return []byte{}, def.KEY_EMPTY
	}
	val, err := s.Client.Get(key)
	if err != nil {
		return []byte{}, err
	}
	if val == nil {
		return []byte{}, nil
	}
	return val, nil

}

//GetBatch retrieves values for given pair of keys in batch
func (s *StoreClient) GetBatch(keys [][]byte) ([][]byte, error) {
	if len(keys) == 0 {
		return [][]byte{}, def.KEY_EMPTY
	}

	val, err := s.Client.BatchGet(keys)
	if err != nil {
		return [][]byte{}, err
	}
	return val, nil
}

//DeleteKey deletes given key from TiKV
func (s *StoreClient) DeleteKey(key []byte) error {
	if len(key) == 0 {
		return def.EMPTY_KEY_CANNOT_BE_DELETED
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
		return def.START_OR_END_KEY_EMPTY
	}
	err := s.Client.DeleteRange(startKey, endKey)
	if err != nil {
		return err
	}
	return nil
}

//Scan iterates from startKey to endKey upto within limit for closed set [startKey,endKey]
func (s *StoreClient) Scan(startKey []byte, endKey []byte, limit int) ([][]byte, [][]byte, error) {
	keys := make([][]byte, 0)
	values := make([][]byte, 0)

	keys, values, err := s.Client.Scan(startKey, endKey, limit)
	if err != nil {
		return [][]byte{}, [][]byte{}, err
	}
	//also include endKey
	val, err := s.Client.Get(endKey)
	if err != nil {
		return [][]byte{}, [][]byte{}, err
	}
	keys = append(keys, endKey)
	values = append(values, val)
	return keys, values, nil
}

//ReverseScan takes startKey, endKey and limit to scan in reverse direction in range [endKey,startKey]
//returns key,value,error
func (s *StoreClient) ReverseScan(startKey []byte, endKey []byte, limit int) ([][]byte, [][]byte, error) {
	keys := make([][]byte, 0)
	values := make([][]byte, 0)
	if len(startKey) == 0 {
		return [][]byte{}, [][]byte{}, def.START_KEY_UNKNOWN
	}
	keys, values, err := s.Client.ReverseScan(startKey, endKey, limit)
	if err != nil {
		return [][]byte{}, [][]byte{}, err
	}
	//also include startKey
	val, err := s.Client.Get(startKey)
	if err != nil {
		return [][]byte{}, [][]byte{}, err
	}
	keys = append(keys, endKey)
	values = append(values, val)
	return keys, values, nil
}

//PrefixScan
//TODO:implement this method
func (s *StoreClient) PrefixScan(startKey []byte, prefix []byte, limit int) ([][]byte, [][]byte, error) {
	return [][]byte{}, [][]byte{}, nil
}

//ReversePrefixScan
//TODO: implement this method
func (s *StoreClient) ReversePrefixScan(endKey []byte, prefix []byte, limit int) ([][]byte, [][]byte, error) {
	return [][]byte{}, [][]byte{}, nil
}
