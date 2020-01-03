package kvstore

import (
	"encoding/binary"
	"errors"
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
func (s *StoreClient) GenerateDBIdentifier() (uint16, error) {
	val, err := s.Get([]byte(META_DBIDENTIFIER))
	if err != nil {
		return uint16(0), err
	}
	if len(val) == 0 {
		identifier := make([]byte, 2)
		binary.LittleEndian.PutUint16(identifier, DBIDENTIFIER_INITIALCOUNT)
		err := s.Put([]byte(META_DBIDENTIFIER), identifier)
		if err != nil {
			return uint16(0), err
		}
		return DBIDENTIFIER_INITIALCOUNT, nil
	} else {
		identifier := binary.LittleEndian.Uint16(val)
		binary.LittleEndian.PutUint16(val, uint16(identifier+1))
		return identifier, nil
	}
}

//GetDBIdentifier returns identifier for given db
func (s *StoreClient) GetDBIdentifier(dbname []byte) (uint16, error) {
	if len(dbname) == 0 {
		return uint16(0), errors.New("dbname empty")
	}
	val, err := s.Get([]byte(META_DB + string(dbname)))
	if err != nil {
		return uint16(0), err
	}
	//if len(val) is zero, generate a new identifier
	if len(val) == 0 {
		identifier, err := s.GenerateDBIdentifier()
		if err != nil {
			return uint16(0), err
		}
		return identifier, nil
	} else {
		identifier := binary.LittleEndian.Uint16(val)
		return identifier, nil
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
	return string(val), nil
}

//GenerateCollectionIdentifier generate collection identifier and increses identifier by 1
func (s *StoreClient) GenerateCollectionIdentifier() (uint32, error) {
	val, err := s.Get([]byte(META_COLLECTIONIDENTIFIER))
	if err != nil {
		return uint32(0), err
	}
	if len(val) == 0 {
		identifier := make([]byte, 4)
		binary.LittleEndian.PutUint32(identifier, COLLECTIONIDENTIFIER_INITIALCOUNT)
		err := s.Put([]byte(META_COLLECTIONIDENTIFIER), identifier)
		if err != nil {
			return uint32(0), err
		}
		return COLLECTIONIDENTIFIER_INITIALCOUNT, nil
	} else {
		identifier := binary.LittleEndian.Uint32(val)
		binary.LittleEndian.PutUint32(val, uint32(identifier+1))
		return identifier, nil
	}
}

//GetCollectionIdentifier returns identifier for given collection
func (s *StoreClient) GetCollectionIdentifier(collection []byte) (uint32, error) {
	if len(collection) == 0 {
		return uint32(0), errors.New("collection name empty")
	}
	val, err := s.Get([]byte(META_COLLECTION + string(collection)))
	if err != nil {
		return uint32(0), err
	}
	//if len(val) is zero, generate a new identifier
	if len(val) == 0 {
		identifier, err := s.GenerateCollectionIdentifier()
		if err != nil {
			return uint32(0), err
		}
		return identifier, nil
	} else {
		identifier := binary.LittleEndian.Uint32(val)
		return identifier, nil
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
func (s *StoreClient) GenerateNamespaceIdentifier() (uint32, error) {
	val, err := s.Get([]byte(META_NAMESPACEIDENTIFIER))
	if err != nil {
		return uint32(0), err
	}
	//if there is no namespace id, generate a new one
	//TODO: move this logic to separate init file for performance
	if len(val) == 0 {
		identifier := make([]byte, 4)
		binary.LittleEndian.PutUint32(identifier, NAMESPACEIDENTIFIER_INITIALCOUNT)
		err := s.Put([]byte(META_NAMESPACEIDENTIFIER), identifier)
		if err != nil {
			return uint32(0), err
		}
		return NAMESPACEIDENTIFIER_INITIALCOUNT, nil
	} else {
		identifier := binary.LittleEndian.Uint32(val)
		binary.LittleEndian.PutUint32(val, uint32(identifier+1))
		return identifier, nil
	}
}

//GetNamespaceIdentifier returns identifier for given namespace
func (s *StoreClient) GetNamespaceIdentifier(namespace []byte) (uint32, error) {
	if len(namespace) == 0 {
		return uint32(0), errors.New("collection name empty")
	}
	val, err := s.Get([]byte(META_NAMESPACE + string(namespace)))
	if err != nil {
		return uint32(0), err
	}
	//if len(val) is zero, generate a new identifier
	if len(val) == 0 {
		identifier, err := s.GenerateNamespaceIdentifier()
		if err != nil {
			return uint32(0), err
		}
		return identifier, nil
	} else { //else send the value read from db
		identifier := binary.LittleEndian.Uint32(val)
		return identifier, nil
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

//GenerateUniqueID generates a 13byte unique_id which composes of
//4 byte UNIX timestamp
//3-byte MAC Addr
//2 byte processID
//4 byte RandomCounter Value
func (s *StoreClient) GenerateUniqueID() []byte {
	unixTimeStamp := getUnixTimestamp() //returns 4 byte UNIX timestamp
	macAddr := getMACAddress()          //returns 3 byte MAC Address
	processID := getProcessID()         //returns 2 byte ProcessID
	counter := generateRandomCount()    //returns 4 byte RANDOM count
	uniqueID := append(unixTimeStamp, macAddr...)
	uniqueID = append(uniqueID, processID...)
	uniqueID = append(uniqueID, counter...)

	return uniqueID
}
