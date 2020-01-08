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
		indexKey := []byte(INDEX_KEY + string(dbID) + ":" + string(collectionID) + ":" + string(namespaceID) + ":" + fieldToIndex + ":" + typeOfData[fieldToIndex] + ":" + string(fieldValue))

		fmt.Println("indexkey: ", (INDEX_KEY + string(dbID) + ":" + string(collectionID) + ":" + string(namespaceID) + ":" + fieldToIndex + ":" + typeOfData[fieldToIndex] + ":" + string(fieldValue)))

		//get value for that index key
		val, err := s.Get(indexKey)
		if err != nil {
			return err
		}

		val = append(val, marshaledRB...)
		//add to batch KV pair
		//write in batch
		err = s.Put(indexKey, val)
		if err != nil {
			return err
		}

	}

	return nil
}
