package engine

import (
	"db-arch/server/internal/def"
	"db-arch/server/io"
	"encoding/binary"
	"encoding/json"
	"fmt"

	"github.com/RoaringBitmap/roaring"
)

/*
Design considerations`
---------------------
A typical key consists of following parts:

- db_name [2 bytes] ~ 65k values
- collection_name [4 bytes]
- namespace [4 bytes]
- unique_id [4 bytes]
Total key size for a document will be 14 bytes.
*/

type Engine struct {
	DBName      string
	DBID        []byte
	Namespace   string
	NamespaceID []byte
	Session     map[string][]byte //session is used to check whether given d,c,n creds are correct
}

//TODO: verify
//ConnectDB initializes engine with DBName, Namespace
func (e *Engine) ConnectDB(s io.Store) error {
	fmt.Println("[[ConnectDB]] inside function")
	dbname := []byte(e.DBName)
	namespace := []byte(e.Namespace)

	dbID, err := e.GetDBIdentifier(s, dbname)
	if err != nil {
		return err
	}

	e.Session = make(map[string][]byte)

	e.Session[e.DBName] = dbID
	e.DBID = dbID

	namespaceID, err := e.GetNamespaceIdentifier(s, namespace)
	if err != nil {
		return err
	}

	e.Session[e.Namespace] = namespaceID
	e.NamespaceID = namespaceID

	return nil
}

//TODO: verify
//GenerateDBIdentifier return db identifier value and increase identifier by 1
func (e *Engine) GenerateDBIdentifier(s io.Store, dbname []byte) ([]byte, error) {
	val, err := s.Get([]byte(def.META_DBIDENTIFIER))
	if err != nil {
		return []byte{}, err
	}
	//if there is no id present, generate a new one
	if len(val) == 0 {
		identifier := make([]byte, 2)

		binary.BigEndian.PutUint16(identifier, def.DBIDENTIFIER_INITIALCOUNT)
		err := s.Put([]byte(def.META_DBIDENTIFIER), identifier)
		if err != nil {
			return []byte{}, err
		}

		return identifier, nil
	} else {
		identifier := binary.BigEndian.Uint16(val)
		binary.BigEndian.PutUint16(val, uint16(identifier+1))

		err := s.Put([]byte(def.META_DBIDENTIFIER), val)
		if err != nil {
			return []byte{}, err
		}
		return val, nil
	}
}

//TODO: verify
//GetDBIdentifier returns identifier for given db
func (e *Engine) GetDBIdentifier(s io.Store, dbname []byte) ([]byte, error) {
	if len(dbname) == 0 {
		return []byte{}, def.DB_NAME_EMPTY
	}
	val, err := s.Get([]byte(def.META_DB + string(dbname)))
	if err != nil {
		return []byte{}, err
	}

	//if len(val) is zero, generate a new identifier
	if len(val) == 0 {
		identifier, err := e.GenerateDBIdentifier(s, dbname)
		if err != nil {
			return []byte{}, err
		}

		//insert meta:db:dname = identifier
		err = s.Put([]byte(def.META_DB+string(dbname)), identifier)
		if err != nil {
			return []byte{}, err
		}

		//insert meta:dbid:id=name
		err = s.Put(append([]byte(def.META_DBID), identifier...), dbname)
		if err != nil {
			return []byte{}, err
		}

		return identifier, nil
	} else {
		//identifier := binary.
		//
		//LittleEndian.Uint16(val)
		return val, nil
	}
}

//TODO: verify
//GetDBName returns database name for given db identifier
func (e *Engine) GetDBName(s io.Store, dbIdentifier []byte) (string, error) {
	if len(dbIdentifier) == 0 {
		return "", def.DB_IDENTIFIER_EMPTY
	}
	val, err := s.Get(append([]byte(def.META_DBID), dbIdentifier...))
	if err != nil {
		return "", err
	}
	if val == nil {
		return "", nil
	}
	return string(val), nil
}

//TODO: verify
//GenerateCollectionIdentifier generate collection identifier and increases identifier by 1
func (e *Engine) GenerateCollectionIdentifier(s io.Store, collectionname []byte) ([]byte, error) {
	val, err := s.Get([]byte(def.META_COLLECTIONIDENTIFIER))
	if err != nil {
		return []byte{}, err
	}
	if len(val) == 0 {
		identifier := make([]byte, 4)
		binary.BigEndian.PutUint32(identifier, def.COLLECTIONIDENTIFIER_INITIALCOUNT)
		err := s.Put([]byte(def.META_COLLECTIONIDENTIFIER), identifier)
		if err != nil {
			return []byte{}, err
		}
		return identifier, nil
	} else {
		identifier := binary.BigEndian.Uint32(val)
		binary.BigEndian.PutUint32(val, uint32(identifier+1))
		err := s.Put([]byte(def.META_COLLECTIONIDENTIFIER), val)
		if err != nil {
			return []byte{}, err
		}
		return val, nil
	}
}

//TODO: verify
//GetCollectionIdentifier returns identifier for given collection
func (e *Engine) GetCollectionIdentifier(s io.Store, collection []byte) ([]byte, error) {
	if len(collection) == 0 {
		return []byte{}, def.COLLECTION_NAME_EMPTY
	}
	val, err := s.Get([]byte(def.META_COLLECTION + string(collection)))
	if err != nil {
		return []byte{}, err
	}
	//if len(val) is zero, generate a new identifier
	if len(val) == 0 {
		identifier, err := e.GenerateCollectionIdentifier(s, collection)
		if err != nil {
			return []byte{}, err
		}

		//insert meta:collection:collectionname = identifier
		err = s.Put([]byte(def.META_COLLECTION+string(collection)), identifier)
		if err != nil {
			return []byte{}, err
		}

		//insert meta:collectionid:id=name
		err = s.Put(append([]byte(def.META_COLLECTIONID), identifier...), collection)
		if err != nil {
			return []byte{}, err
		}

		return identifier, nil
	} else {
		//identifier := binary.LittleEndian.Uint32(val)
		return val, nil
	}
}

//TODO: verify
//GetCollectionName returns collection name for given collection identifier
func (e *Engine) GetCollectionName(s io.Store, collectionIdentifier []byte) (string, error) {
	if len(collectionIdentifier) == 0 {
		return "", def.COLLECTION_IDENTIFIER_EMPTY
	}
	val, err := s.Get([]byte(def.META_COLLECTIONID + string(collectionIdentifier)))
	if err != nil {
		return "", err
	}
	return string(val), nil
}

//TODO: verify
//GenerateNamespaceIdentifier generates namespace identifier value and increases identifier by 1
func (e *Engine) GenerateNamespaceIdentifier(s io.Store, namespace []byte) ([]byte, error) {
	val, err := s.Get([]byte(def.META_NAMESPACEIDENTIFIER))
	if err != nil {
		return []byte{}, err
	}
	//if there is no namespace id, generate a new one
	//TODO: move this logic to separate init file for performance
	if len(val) == 0 {
		identifier := make([]byte, 4)
		binary.BigEndian.PutUint32(identifier, def.NAMESPACEIDENTIFIER_INITIALCOUNT)
		err := s.Put([]byte(def.META_NAMESPACEIDENTIFIER), identifier)
		if err != nil {
			return []byte{}, err
		}
		return identifier, nil
	} else {
		identifier := binary.BigEndian.Uint32(val)
		binary.BigEndian.PutUint32(val, uint32(identifier+1))
		err := s.Put([]byte(def.META_NAMESPACEIDENTIFIER), val)
		if err != nil {
			return []byte{}, err
		}
		return val, nil
	}
}

//TODO: verify
//GetNamespaceIdentifier returns identifier for given namespace
func (e *Engine) GetNamespaceIdentifier(s io.Store, namespace []byte) ([]byte, error) {
	if len(namespace) == 0 {
		return []byte{}, def.COLLECTION_NAME_EMPTY
	}
	val, err := s.Get([]byte(def.META_NAMESPACE + string(namespace)))
	if err != nil {
		return []byte{}, err
	}
	//if len(val) is zero, generate a new identifier
	if len(val) == 0 {
		identifier, err := e.GenerateNamespaceIdentifier(s, namespace)
		if err != nil {
			return []byte{}, err
		}

		//insert meta:namespace:namespace = identifier
		err = s.Put([]byte(def.META_NAMESPACE+string(namespace)), identifier)
		if err != nil {
			return []byte{}, err
		}

		//insert meta:namespaceid:id=name
		err = s.Put(append([]byte(def.META_NAMESPACEID), identifier...), namespace)
		if err != nil {
			return []byte{}, err
		}

		return identifier, nil
	} else { //else send the value read from db
		//identifier := binary.LittleEndian.Uint32(val)
		return val, nil
	}
}

//TODO: verify
//GetNamespaceName return namespace name with given namespace identifier
func (e *Engine) GetNamespaceName(s io.Store, namespaceIdentifier []byte) (string, error) {
	if len(namespaceIdentifier) == 0 {
		return "", def.NAMESPACE_IDENTIFIER_EMPTY
	}
	val, err := s.Get([]byte(def.META_NAMESPACEID + string(namespaceIdentifier)))
	if err != nil {
		return "", err
	}
	return string(val), err
}

//TODO: verify
//GenerateUniqueID generates a 4 byte unique_id
func (e *Engine) GenerateUniqueID(s io.Store, collectionID []byte) ([]byte, error) {
	//unixTimeStamp := getUnixTimestamp() //returns 4 byte UNIX timestamp
	//macAddr := getMACAddress()          //returns 3 byte MAC Address
	//processID := getProcessID()         //returns 2 byte ProcessID
	//counter := generateRandomCount()    //returns 4 byte RANDOM count
	//uniqueID := append(unixTimeStamp, macAddr...)
	//uniqueID = append(uniqueID, processID...)
	//uniqueID = append(uniqueID, counter...)

	//key format: _uniqueid:dbid:colid:namespaceid=idcounter
	idKey := []byte(def.UNIQUE_ID + string(e.DBID) + ":" + string(collectionID) + ":" + string(e.NamespaceID))
	idCounterInBytes, err := s.Get(idKey)
	if err != nil {
		return []byte{}, err
	}
	//if there is no id counter, generate a new one
	if len(idCounterInBytes) == 0 {
		//create 4 byte identifier for each document
		counterByte := make([]byte, 4)
		binary.BigEndian.PutUint32(counterByte, def.UNIQUE_ID_INITIALCOUNT)
		err := s.Put(idKey, counterByte)
		if err != nil {
			return []byte{}, err
		}
		return counterByte, nil
	} else {
		currentCount := binary.BigEndian.Uint32(idCounterInBytes)
		counterByte := make([]byte, 4)
		//increase count by 1 and write to db
		binary.BigEndian.PutUint32(counterByte, (currentCount + 1))
		//insert
		err := s.Put(idKey, counterByte)
		if err != nil {
			return []byte{}, err
		}
		return counterByte, nil
	}
}

//TODO: verify
//GetIdentifiers returns database, collection and namespace identifiers for respective names given
//and generate new ones if they do not exist
func (e *Engine) GetIdentifiers(s io.Store, database string, collection string,
	namespace string) ([]byte, []byte, []byte, error) {

	dbID, err := e.GetDBIdentifier(s, []byte(database))
	if err != nil {
		return []byte{}, []byte{}, []byte{}, err
	}

	collectionID, err := e.GetCollectionIdentifier(s, []byte(collection))
	if err != nil {
		return []byte{}, []byte{}, []byte{}, err
	}

	namespaceID, err := e.GetNamespaceIdentifier(s, []byte(namespace))
	if err != nil {
		return []byte{}, []byte{}, []byte{}, err
	}

	return dbID, collectionID, namespaceID, nil
}

//TODO: verify
//SearchIdentifiers retrieves db,collection,namespace identifiers only if they exist
func (e *Engine) SearchIdentifiers(s io.Store, dbname string, collection string,
	namespace string) ([]byte, []byte, []byte, error) {

	dbID, err := s.Get([]byte(def.META_DB + string(dbname)))
	if len(dbID) == 0 || err != nil {
		return []byte{}, []byte{}, []byte{}, err
	}

	collectionID, err := s.Get([]byte(def.META_COLLECTION + string(collection)))
	if len(collectionID) == 0 || err != nil {
		return []byte{}, []byte{}, []byte{}, err
	}

	namespaceID, err := s.Get([]byte(def.META_NAMESPACE + string(namespace)))
	if len(namespaceID) == 0 || err != nil {
		return []byte{}, []byte{}, []byte{}, err
	}

	return dbID, collectionID, namespaceID, nil
}

//TODO: verify
//InsertDocument retrieves identifiers and inserts document to database

func (e *Engine) InsertDocument(s io.Store,
	collection string,
	data map[string][]byte, indices []string) error {

	if len(e.DBName) == 0 || len(collection) == 0 || len(e.Namespace) == 0 {
		return def.NAMES_CANNOT_BE_EMPTY
	}

	//KV pair to insert in batch
	keyCache := make([][]byte, 0)
	valueCache := make([][]byte, 0)

	//get database, collection,namespace identifiers
	//dbID, collectionID, namespaceID, err := e.GetIdentifiers(s, database, collection, namespace)
	//if err != nil {
	//urn err
	//}

	if _, ok := e.Session[e.DBName]; !ok {
		return def.DB_DOES_NOT_EXIST
	}

	if _, ok := e.Session[e.Namespace]; !ok {
		return def.NAMESPACE_DOES_NOT_EXIST
	}

	collectionID, err := e.GetCollectionIdentifier(s, []byte(collection))

	if err != nil {
		return err
	}

	//generate unique_id
	uniqueID, err := e.GenerateUniqueID(s, collectionID)
	if err != nil {
		return err
	}
	//generate key
	fmt.Println("[[engine.go/uniqueID]]", uniqueID)
	key := []byte(string(e.DBID) + ":" + string(collectionID) + ":" + string(e.NamespaceID) + ":" + string(uniqueID))

	dataInBytes, err := json.Marshal(data)
	if err != nil {
		return err
	}

	//insert into db
	//err = s.Put(key, dataInBytes)
	//if err != nil {
	//	return err
	//}
	keyCache = append(keyCache, key)
	valueCache = append(valueCache, dataInBytes)
	//indexer
	indexKey, indexValue, err := e.IndexDocument(s, collectionID, uniqueID, data, indices)
	if err != nil {
		return err
	}

	keyCache = append(keyCache, indexKey...)
	valueCache = append(valueCache, indexValue...)

	//insert in batch
	err = s.PutBatch(keyCache, valueCache)
	if err != nil {
		return err
	}

	return nil
}

//TODO:  verify
//SearchDocument queries document for given query params
func (e *Engine) SearchDocument(s io.Store, collection string,
	query []string) ([][]byte, error) {

	if len(e.DBName) == 0 || len(collection) == 0 || len(e.Namespace) == 0 {
		return [][]byte{}, def.NAMES_CANNOT_BE_EMPTY
	}

	//get identifiers for given
	//dbID, collectionID, namespaceID, err := e.SearchIdentifiers(s, database, collection, namespace)
	//if err != nil {
	//	return [][]byte{}, err
	//}

	if _, ok := e.Session[e.DBName]; !ok {
		return [][]byte{}, def.DB_DOES_NOT_EXIST
	}

	if _, ok := e.Session[e.Namespace]; !ok {
		return [][]byte{}, def.NAMESPACE_DOES_NOT_EXIST
	}

	//here if collection doesn't exist, do not create new one
	collectionID, err := s.Get([]byte(def.META_COLLECTION + collection))
	if err != nil {
		return [][]byte{}, err
	}

	//collectionID check is required here
	if len(e.DBID) == 0 || len(collectionID) == 0 || len(e.DBID) == 0 {
		return [][]byte{}, def.IDENTIFIER_NOT_FOUND
	}

	/*
		// //find typeOfData  and get byteOrderedData
		// typeOfData, byteOrderedData := findTypeOfData(query)

		// rb := roaring.New()

		// for fieldName, fieldType := range typeOfData {

		// 	//generate indexKey
		// 	indexKey := []byte(def.INDEX_KEY + string(dbID) + ":" + string(collectionID) + ":" + string(namespaceID) + ":" + fieldName + ":" + fieldType + ":" + string(byteOrderedData[fieldName]))

		// 	uniqueIDBitmapArray, err := s.Get(indexKey)
		// 	if len(uniqueIDBitmapArray) == 0 || err != nil {
		// 		return [][]byte{}, err
		// 	}

		// 	if rb.IsEmpty() == true {

		// 		err := rb.UnmarshalBinary(uniqueIDBitmapArray)

		// 		if err != nil {
		// 			return [][]byte{}, nil

		// 		}
		// 	} else {

		// 		tmp := roaring.New()
		// 		err := tmp.UnmarshalBinary(uniqueIDBitmapArray)
		// 		if err != nil {
		// 			return [][]byte{}, err
		// 		}
		// 		rb = roaring.FastAnd(rb, tmp) //fast AND two bitmaps
		// 	}

		// }

		// if rb.IsEmpty() == true {
		// 	return [][]byte{}, nil
		// }

	*/

	fmt.Println("[[engine.go]] evaluate postfix expression")
	rb, err := e.EvaluatePostFix(s, query, collectionID)
	if err != nil {
		return [][]byte{}, err
	}

	resultRoaring := rb.(roaring.Bitmap)
	//fmt.Println("FINAL RB is ", resultRoaring.ToArray())
	//retrieve document keys for search
	fmt.Println("result roaring size -->", resultRoaring.GetSerializedSizeInBytes())
	fmt.Println("length of roaring -->", len(resultRoaring.ToArray()))
	searchKeys := make([][]byte, 0)
	searchKeyLength := len(resultRoaring.ToArray())
	uniqueIDArr := resultRoaring.ToArray() //get all IDs
	//get all documents keys
	for i := 0; i < searchKeyLength; i++ {
		uniqueIDByte := make([]byte, 4)

		binary.BigEndian.PutUint32(uniqueIDByte, uniqueIDArr[i])
		documentKeys := []byte(string(e.DBID) + ":" + string(collectionID) + ":" + string(e.NamespaceID) + ":" + string(uniqueIDByte))
		searchKeys = append(searchKeys, documentKeys)
	}
	resultArr, err := s.GetBatch(searchKeys)
	if err != nil {
		return [][]byte{}, err
	}
	return resultArr, nil
}
