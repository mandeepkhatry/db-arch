package kvstore

import (
	"encoding/binary"

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

	batchKV := make([]byte, 0)

	//convert uniqueID into uint32 and change into roaring bitmap
	num := binary.LittleEndian.Uint32(uniqueID)
	rb := roaring.BitmapOf(num)
	marshaledRB, err := rb.MarshalBinary()
	if err != nil {
		return err
	}

	for i := 0; i < len(indices); i++ {
		fieldToIndex := indices[i]
		if typeOfData[fieldToIndex] == "words" {
			//TODO: implementation for strings yet to be done
			//val := newData[fieldToIndex]
			//splitWords := strings.Split(string(val), " ")
			//for k, v := range splitWords {
			//	indexKey := generateKey([]byte(INDEX_KEY), dbID, collectionID, namespaceID,
			//		[]byte(fieldToIndex), []byte(typeOfData["fieldToIndex"]), []byte(v))
			//	s.Get()
			//}
		} else {
			fieldValue := newData[fieldToIndex]

			//generate index key
			indexKey := generateKey([]byte(INDEX_KEY), dbID, collectionID, namespaceID,
				[]byte(fieldToIndex), []byte(typeOfData[fieldToIndex]), []byte(fieldValue))

			//get value for that index key
			val, err := s.Get(indexKey)
			if err != nil {
				return err
			}

			val = append(val, marshaledRB...)
			//add to batch KV pair
			batchKV = append(batchKV, indexKey...)
			batchKV = append(batchKV, val...)

		}

	}

	//write in batch
	err = s.PutBatch(batchKV)
	if err != nil {
		return err
	}
	return nil
}
