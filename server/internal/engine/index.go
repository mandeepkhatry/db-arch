package engine

import (
	"db-arch/server/internal/def"
	"db-arch/server/io"
	"encoding/binary"
	"fmt"
	"strings"

	"github.com/RoaringBitmap/roaring"
)

//TODO: verify
//IndexDocument indexes document in batch
func (e *Engine) IndexDocument(s io.Store, collectionID []byte,
	uniqueID []byte, data map[string][]byte, indices []string) ([][]byte, [][]byte, error) {

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
	fmt.Println("TYPE OF DATA is ", typeOfData)
	//convert uniqueID into uint32
	num := binary.BigEndian.Uint32(uniqueID)
	fmt.Println("[[index.go]]uniqueID in int32:", num)
	arrKeys := make([][]byte, 0)
	arrValues := make([][]byte, 0)

	for i := 0; i < len(indices); i++ {
		//remove all whitespaces if any
		indexStr := strings.Replace(indices[i], " ", "", -1)
		indexStrArr := strings.Split(indexStr, ",")

		//if there is only one condition
		if len(indexStrArr) == 1 {

			fieldToIndex := indices[i]
			//TODO: tokenize words and create index for them too

			fieldValue := newData[fieldToIndex]

			//generate index key
			indexKey := []byte(def.INDEX_KEY + string(e.DBID) + ":" + string(collectionID) + ":" + string(e.NamespaceID) + ":" + fieldToIndex + ":" + typeOfData[fieldToIndex] + ":" + string(fieldValue))
			fmt.Println("NEW INDEXKEY : ", indexKey)

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
				//err = s.Put(indexKey, marshaledRB)
				if err != nil {
					return [][]byte{}, [][]byte{}, err
				}
				arrKeys = append(arrKeys, indexKey)
				arrValues = append(arrValues, marshaledRB)
			} else {

				rb := roaring.BitmapOf(num)
				marshaledRB, err := rb.MarshalBinary()
				if err != nil {
					return [][]byte{}, [][]byte{}, err
				}

				arrKeys = append(arrKeys, indexKey)
				arrValues = append(arrValues, marshaledRB)

			}

		} else {
			//for compound index
			/*
				parse fieldname,type for each compound condition and create compound index
			*/

			indexKey := def.INDEX_KEY + string(e.DBID) + ":" + string(collectionID) + ":" + string(e.NamespaceID)

			for _, fieldToIndex := range indexStrArr {
				fieldValue := newData[fieldToIndex]
				indexKey += ":" + fieldToIndex + ":" + typeOfData[fieldToIndex] + ":" + string(fieldValue)
			}
			//retrieve set of IDs for each single key and perform AND operation on them
			//singleIndexKey := []byte(def.INDEX_KEY + string(e.DBID) + ":" + string(collectionID) + ":" + string(e.NamespaceID) + ":" + fieldToIndex + ":" + typeOfData[fieldToIndex] + ":" + string(fieldValue))

			val, err := s.Get([]byte(indexKey))
			if err != nil {
				return [][]byte{}, [][]byte{}, err
			}

			if len(val) != 0 {
				tmp := roaring.New()
				err := tmp.UnmarshalBinary(val)
				if err != nil {
					return [][]byte{}, [][]byte{}, err
				}
				tmpArr := tmp.ToArray()
				tmpArr = append(tmpArr, num)

				rb := roaring.BitmapOf(tmpArr...)
				marshaledRB, err := rb.MarshalBinary()
				if err != nil {
					return [][]byte{}, [][]byte{}, err
				}

				arrKeys = append(arrKeys, []byte(indexKey))
				arrValues = append(arrValues, marshaledRB)

			} else {
				rb := roaring.BitmapOf(num)
				marshaledRB, err := rb.MarshalBinary()
				if err != nil {
					return [][]byte{}, [][]byte{}, err
				}

				arrKeys = append(arrKeys, []byte(indexKey))
				arrValues = append(arrValues, marshaledRB)

			}
		}

	}

	fmt.Println("[[index.go]]arrKeys:", arrKeys)
	fmt.Println("[[index.go]]arrValues:", arrValues)
	return arrKeys, arrValues, nil
}
