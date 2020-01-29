package engine

import (
	"db-arch/server/internal/def"
	"db-arch/server/io"
	"encoding/binary"

	"github.com/RoaringBitmap/roaring"
)

//TODO: verify
//IndexDocument indexes document in batch
func (e *Engine) IndexDocument(s io.Store, dbID []byte, collectionID []byte,
	namespaceID []byte, uniqueID []byte,
	data map[string][]byte, indices []string) ([][]byte, [][]byte, error) {

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

	//convert uniqueID into uint32
	num := binary.LittleEndian.Uint32(uniqueID)
	arrKeys := make([][]byte, 0)
	arrValues := make([][]byte, 0)

	for i := 0; i < len(indices); i++ {
		fieldToIndex := indices[i]
		//TODO: tokenize words and create index for them too

		fieldValue := newData[fieldToIndex]

		//generate index key
		indexKey := []byte(def.INDEX_KEY + string(dbID) + ":" + string(collectionID) + ":" + string(namespaceID) + ":" + fieldToIndex + ":" + typeOfData[fieldToIndex] + ":" + string(fieldValue))

		//get value for that index key
		val, err := s.Get(indexKey)
		if err != nil {
			return [][]byte{}, [][]byte{}, err
		}
		//if index already exists, append uniqueIDs
		if len(val) != 0 {
			tmp := roaring.New()
			err = tmp.UnmarshalBinary(val)
			if err != nil {
				return [][]byte{}, [][]byte{}, err
			}
			tmpArr := tmp.ToArray()
			tmpArr = append(tmpArr, num)

			rb := roaring.BitmapOf(tmpArr...)
			marshaledRB, err := rb.MarshalBinary()
			//add to DB
			err = s.Put(indexKey, marshaledRB)
			if err != nil {
				return [][]byte{}, [][]byte{}, err
			}
		} else {

			rb := roaring.BitmapOf(num)
			marshaledRB, err := rb.MarshalBinary()
			if err != nil {
				return [][]byte{}, [][]byte{}, err
			}

			arrKeys = append(arrKeys, indexKey)
			arrValues = append(arrValues, marshaledRB)

		}

	}
	return arrKeys, arrValues, nil
}
