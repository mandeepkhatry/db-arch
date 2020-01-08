package kvstore

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/RoaringBitmap/roaring"
)

/*
Design considerations
---------------------
A typical key consists of following parts:

- db_name [2 bytes] ~ 65k values
- collection_name [4 bytes]
- namespace [4 bytes]
- unique_id [13 byte]
	(unique_id is essential component of key. We took some favor of mongodb way of making unique id.)
	- unix timestamp [4 byte]
	- MAC address of current machine [3 byte]
	- process_id [2 byte]
	- a 4-byte counter starting with a random value with unix timestamp as seed

Total key size for a document will be 20 bytes.
*/

//GenerateDBIdentifier return db identifier value and increase identifier by 1
func (s *StoreClient) GenerateDBIdentifier(dbname []byte) ([]byte, error) {
	val, err := s.Get([]byte(META_DBIDENTIFIER))
	if err != nil {
		return []byte{}, err
	}

	if len(val) == 0 {
		identifier := make([]byte, 2)

		binary.LittleEndian.PutUint16(identifier, DBIDENTIFIER_INITIALCOUNT)
		err := s.Put([]byte(META_DBIDENTIFIER), identifier)
		if err != nil {
			return []byte{}, err
		}

		return identifier, nil
	} else {
		identifier := binary.LittleEndian.Uint16(val)
		binary.LittleEndian.PutUint16(val, uint16(identifier+1))

		err := s.Put([]byte(META_DBIDENTIFIER), val)
		if err != nil {
			return []byte{}, err
		}
		return val, nil
	}
}

//GetDBIdentifier returns identifier for given db
func (s *StoreClient) GetDBIdentifier(dbname []byte) ([]byte, error) {
	if len(dbname) == 0 {
		return []byte{}, errors.New("dbname empty")
	}
	val, err := s.Get([]byte(META_DB + string(dbname)))
	fmt.Println("[[GetDBIdentifier]] value: ", string(val))
	if err != nil {
		return []byte{}, err
	}
	//if len(val) is zero, generate a new identifier
	if len(val) == 0 {
		identifier, err := s.GenerateDBIdentifier(dbname)
		fmt.Println("[[GetDBIdentifier]] identifier", string(identifier))
		if err != nil {
			return []byte{}, err
		}

		//insert meta:db:dname = identifier
		err = s.Put([]byte(META_DB+string(dbname)), identifier)
		if err != nil {
			return []byte{}, err
		}

		//insert meta:dbid:id=name
		err = s.Put(append([]byte(META_DBID), identifier...), dbname)
		if err != nil {
			return []byte{}, err
		}

		return identifier, nil
	} else {
		//identifier := binary.LittleEndian.Uint16(val)
		return val, nil
	}
}

//GetDBName returns database name for given db identifier
func (s *StoreClient) GetDBName(dbIdentifier []byte) (string, error) {
	if len(dbIdentifier) == 0 {
		return "", errors.New("dbidentifier empty")
	}
	val, err := s.Get(append([]byte(META_DBID), dbIdentifier...))
	if err != nil {
		return "", err
	}
	if val == nil {
		return "", nil
	}
	return string(val), nil
}

//GenerateCollectionIdentifier generate collection identifier and increases identifier by 1
func (s *StoreClient) GenerateCollectionIdentifier(collectionname []byte) ([]byte, error) {
	val, err := s.Get([]byte(META_COLLECTIONIDENTIFIER))
	if err != nil {
		return []byte{}, err
	}
	if len(val) == 0 {
		identifier := make([]byte, 4)
		binary.LittleEndian.PutUint32(identifier, COLLECTIONIDENTIFIER_INITIALCOUNT)
		err := s.Put([]byte(META_COLLECTIONIDENTIFIER), identifier)
		if err != nil {
			return []byte{}, err
		}
		return identifier, nil
	} else {
		identifier := binary.LittleEndian.Uint32(val)
		binary.LittleEndian.PutUint32(val, uint32(identifier+1))
		err := s.Put([]byte(META_COLLECTIONIDENTIFIER), val)
		if err != nil {
			return []byte{}, err
		}
		return val, nil
	}
}

//GetCollectionIdentifier returns identifier for given collection
func (s *StoreClient) GetCollectionIdentifier(collection []byte) ([]byte, error) {
	if len(collection) == 0 {
		return []byte{}, errors.New("collection name empty")
	}
	val, err := s.Get([]byte(META_COLLECTION + string(collection)))
	if err != nil {
		return []byte{}, err
	}
	//if len(val) is zero, generate a new identifier
	if len(val) == 0 {
		identifier, err := s.GenerateCollectionIdentifier(collection)
		if err != nil {
			return []byte{}, err
		}

		//insert meta:collection:collectionname = identifier
		err = s.Put([]byte(META_COLLECTION+string(collection)), identifier)
		if err != nil {
			return []byte{}, err
		}

		//insert meta:collectionid:id=name
		err = s.Put(append([]byte(META_COLLECTIONID), identifier...), collection)
		if err != nil {
			return []byte{}, err
		}

		return identifier, nil
	} else {
		//identifier := binary.LittleEndian.Uint32(val)
		return val, nil
	}
}

//GetCollectionName returns collection name for given collection identifier
func (s *StoreClient) GetCollectionName(collectionIdentifier []byte) (string, error) {
	if len(collectionIdentifier) == 0 {
		return "", errors.New("collection identifier empty")
	}
	val, err := s.Get([]byte(META_COLLECTIONID + string(collectionIdentifier)))
	if err != nil {
		return "", err
	}
	return string(val), nil
}

//GenerateNamespaceIdentifier generates namespace identifier value and increases identifier by 1
func (s *StoreClient) GenerateNamespaceIdentifier(namespace []byte) ([]byte, error) {
	val, err := s.Get([]byte(META_NAMESPACEIDENTIFIER))
	if err != nil {
		return []byte{}, err
	}
	//if there is no namespace id, generate a new one
	//TODO: move this logic to separate init file for performance
	if len(val) == 0 {
		identifier := make([]byte, 4)
		binary.LittleEndian.PutUint32(identifier, NAMESPACEIDENTIFIER_INITIALCOUNT)
		err := s.Put([]byte(META_NAMESPACEIDENTIFIER), identifier)
		if err != nil {
			return []byte{}, err
		}
		return identifier, nil
	} else {
		identifier := binary.LittleEndian.Uint32(val)
		binary.LittleEndian.PutUint32(val, uint32(identifier+1))
		err := s.Put([]byte(META_NAMESPACEIDENTIFIER), val)
		if err != nil {
			return []byte{}, err
		}
		return val, nil
	}
}

//GetNamespaceIdentifier returns identifier for given namespace
func (s *StoreClient) GetNamespaceIdentifier(namespace []byte) ([]byte, error) {
	if len(namespace) == 0 {
		return []byte{}, errors.New("collection name empty")
	}
	val, err := s.Get([]byte(META_NAMESPACE + string(namespace)))
	if err != nil {
		return []byte{}, err
	}
	//if len(val) is zero, generate a new identifier
	if len(val) == 0 {
		identifier, err := s.GenerateNamespaceIdentifier(namespace)
		if err != nil {
			return []byte{}, err
		}

		//insert meta:namespace:namespace = identifier
		err = s.Put([]byte(META_NAMESPACE+string(namespace)), identifier)
		if err != nil {
			return []byte{}, err
		}

		//insert meta:namespaceid:id=name
		err = s.Put(append([]byte(META_NAMESPACEID), identifier...), namespace)
		if err != nil {
			return []byte{}, err
		}

		return identifier, nil
	} else { //else send the value read from db
		//identifier := binary.LittleEndian.Uint32(val)
		return val, nil
	}
}

//GetNamespaceName return namespace name with given namespace identifier
func (s *StoreClient) GetNamespaceName(namespaceIdentifier []byte) (string, error) {
	if len(namespaceIdentifier) == 0 {
		return "", errors.New("namespace identifier empty")
	}
	val, err := s.Get([]byte(META_NAMESPACEID + string(namespaceIdentifier)))
	if err != nil {
		return "", err
	}
	return string(val), err
}

//GenerateUniqueID generates a 4 byte unique_id
func (s *StoreClient) GenerateUniqueID(dbID []byte, collectionID []byte, namespaceID []byte) ([]byte, error) {
	//unixTimeStamp := getUnixTimestamp() //returns 4 byte UNIX timestamp
	//macAddr := getMACAddress()          //returns 3 byte MAC Address
	//processID := getProcessID()         //returns 2 byte ProcessID
	//counter := generateRandomCount()    //returns 4 byte RANDOM count
	//uniqueID := append(unixTimeStamp, macAddr...)
	//uniqueID = append(uniqueID, processID...)
	//uniqueID = append(uniqueID, counter...)

	//key format: _uniqueid:dbid:colid:namespaceid=idcounter
	idKey := []byte(UNIQUE_ID + string(dbID) + ":" + string(collectionID) + ":" + string(namespaceID))
	idCounterInBytes, err := s.Get(idKey)
	if err != nil {
		return []byte{}, err
	}
	//if there is no id counter, generate a new one
	if len(idCounterInBytes) == 0 {
		//create 4 byte identifier for each document
		counterByte := make([]byte, 4)
		binary.LittleEndian.PutUint32(counterByte, UNIQUE_ID_INITIALCOUNT)
		err := s.Put(idKey, counterByte)
		if err != nil {
			return []byte{}, err
		}
		return counterByte, nil
	} else {
		currentCount := binary.LittleEndian.Uint32(idCounterInBytes)
		counterByte := make([]byte, 4)
		//increase count by 1 and write to db
		binary.LittleEndian.PutUint32(counterByte, (currentCount + 1))
		//insert
		err := s.Put(idKey, counterByte)
		if err != nil {
			return []byte{}, err
		}
		return counterByte, nil
	}
}

//GetIdentifiers returns database, collection and namespace identifiers for respective names given
func (s *StoreClient) GetIdentifiers(database string, collection string,
	namespace string) ([]byte, []byte, []byte, error) {
	dbID, err := s.GetDBIdentifier([]byte(database))
	if err != nil {
		return []byte{}, []byte{}, []byte{}, err
	}

	collectionID, err := s.GetCollectionIdentifier([]byte(collection))
	if err != nil {
		return []byte{}, []byte{}, []byte{}, err
	}

	namespaceID, err := s.GetNamespaceIdentifier([]byte(namespace))
	if err != nil {
		return []byte{}, []byte{}, []byte{}, err
	}

	return dbID, collectionID, namespaceID, nil
}

//InsertDocument retrieves identifiers and inserts document to database
func (s *StoreClient) InsertDocument(
	database string, collection string, namespace string,
	data map[string][]byte, indices []string) error {

	if len(database) == 0 || len(collection) == 0 || len(namespace) == 0 {
		return errors.New("names can't be empty")
	}

	//get database, collection,namespace identifiers
	dbID, collectionID, namespaceID, err := s.GetIdentifiers(database, collection, namespace)
	if err != nil {
		return err
	}
	fmt.Println("db", dbID)
	fmt.Println("collection ", collectionID)
	fmt.Println("namespace ", namespaceID)

	//generate unique_id
	uniqueID, err := s.GenerateUniqueID(dbID, collectionID, namespaceID)
	if err != nil {
		return err
	}
	//generate key
	key := generateKey(dbID, collectionID, namespaceID, uniqueID)
	dataInBytes, err := json.Marshal(data)
	if err != nil {
		return err
	}

	//insert into db
	err = s.Put(key, dataInBytes)
	if err != nil {
		return err
	}

	//indexer
	err = s.IndexDocument(dbID, collectionID, namespaceID, uniqueID, data, indices)
	if err != nil {
		return err
	}
	return nil
}

//SearchDocument queries document for given query params
func (s *StoreClient) SearchDocument(
	database string, collection string,
	namespace string, query map[string][]byte) ([][]byte, error) {

	fmt.Println("[[engine/searchdocument]] start search!")
	if len(database) == 0 || len(collection) == 0 || len(namespace) == 0 {
		return [][]byte{}, errors.New("names can't be empty")
	}

	//get identifiers for given
	dbID, collectionID, namespaceID, err := s.GetIdentifiers(database, collection, namespace)
	if err != nil {
		return [][]byte{}, err
	}

	fmt.Println("[[engine/searchdocument]] db", dbID)
	fmt.Println("[[engine/searchdocument]] collection ", collectionID)
	fmt.Println("[[engine/searchdocument]] namespace ", namespaceID)

	//find typeOfData  and get byteOrderedData
	typeOfData, byteOrderedData := findTypeOfData(query)

	rb := roaring.New()

	for fieldName, fieldType := range typeOfData {

		//generate indexKey
		indexKey := []byte(INDEX_KEY + string(dbID) + ":" + string(collectionID) + ":" + string(namespaceID) + ":" + fieldName + ":" + fieldType + ":" + string(byteOrderedData[fieldName]))

		uniqueIDBitmapArray, err := s.Get(indexKey)
		if len(uniqueIDBitmapArray) == 0 || err != nil {
			return [][]byte{}, err
		}

		if rb.IsEmpty() == true {

			err := rb.UnmarshalBinary(uniqueIDBitmapArray)

			if err != nil {
				return [][]byte{}, nil

			}
		} else {

			tmp := roaring.New()
			err := tmp.UnmarshalBinary(uniqueIDBitmapArray)
			if err != nil {
				return [][]byte{}, err
			}
			rb = roaring.FastAnd(rb, tmp) //fast AND two bitmaps
		}

	}

	if rb.IsEmpty() == true {
		return [][]byte{}, nil
	}

	//retrieve document keys for search
	searchKeys := make([][]byte, 0)
	searchKeyLength := len(rb.ToArray())
	uniqueIDArr := rb.ToArray() //get all IDs

	//get all documents keys
	for i := 0; i < searchKeyLength; i++ {
		uniqueIDByte := make([]byte, 4)
		binary.LittleEndian.PutUint32(uniqueIDByte, uniqueIDArr[i])
		documentKeys := []byte(string(dbID) + ":" + string(collectionID) + ":" + string(namespaceID) + ":" + string(uniqueIDByte))
		searchKeys = append(searchKeys, documentKeys)
	}

	resultArr, err := s.GetBatch(searchKeys)
	if err != nil {
		return [][]byte{}, err
	}

	return resultArr, nil
}
