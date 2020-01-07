package kvstore

import (
	"encoding/binary"
	"fmt"

	"github.com/RoaringBitmap/roaring"
)

//IndexDocument indexes document in batch
func (s *StoreClient) IndexDocument(dbID []byte, collectionID []byte,
	namespaceID []byte, uniqueID []byte,
	data map[string][]byte, indices []string) error {

	/*

		1. Find type of data
		2. Generate index key
		3. Read db
		4. Append unique_id using roaring bitmap
		5. Write to db in batch

	*/
	typeOfData, newData := findTypeOfData(data)
	/*
		typeOfData:
		map['name']='string'
		map['age']='int'
		map['weight']='double'

		newData:
		map['name']=[]byte,
		map['age']=sorted int []byte
		map['weight']=sorted double []byte
	*/

	//batchKV := make([]byte, 0)

	//convert uniqueID into uint32 and change into roaring bitmap
	num := binary.LittleEndian.Uint32(uniqueID)
	rb := roaring.BitmapOf(num)
	marshaledRB, err := rb.MarshalBinary()
	fmt.Println("marshalled binary: ", marshaledRB)
	if err != nil {
		return err
	}

	for i := 0; i < len(indices); i++ {
		fieldToIndex := indices[i]
		//TODO: tokenize words and create index for them too

		fieldValue := newData[fieldToIndex]

		//generate index key
		indexKey := generateKey([]byte(INDEX_KEY), dbID, collectionID, namespaceID,
			[]byte(fieldToIndex), []byte(typeOfData[fieldToIndex]), fieldValue)

		//get value for that index key
		val, err := s.Get(indexKey)
		if err != nil {
			return err
		}

		val = append(val, marshaledRB...)
		//add to batch KV pair
		//batchKV = append(batchKV, indexKey...)
		//batchKV = append(batchKV, val...)
		//write in batch
		err = s.PutBatch(indexKey, val)
		if err != nil {
			return err
		}

	}

	return nil
}
