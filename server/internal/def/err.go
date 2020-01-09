package def

/*
Package def defines constants, error messages and their status codes
*/

import (
	"errors"
)

//ERROR MESSAGES
var (
	DB_NAME_EMPTY               error = errors.New("dbname empty")
	DB_IDENTIFIER_EMPTY         error = errors.New("dbidentifier empty")
	COLLECTION_NAME_EMPTY       error = errors.New("collection name empty")
	COLLECTION_IDENTIFIER_EMPTY error = errors.New("collection identifier empty")
	NAMESPACE_IDENTIFIER_EMPTY  error = errors.New("namespace identifier empty")
	NAMES_CANNOT_BE_EMPTY       error = errors.New("database/collection/namespace names can't be empty")
	KEY_EMPTY                   error = errors.New("key is empty")
	EMPTY_KEY_CANNOT_BE_DELETED error = errors.New("can't delete empty key")
	START_OR_END_KEY_EMPTY      error = errors.New("start or end key is empty")
	START_KEY_UNKNOWN           error = errors.New("Can't scan from last without knowing startKey")
	IDENTIFIER_NOT_FOUND        error = errors.New("id not found for given db/collection/namespace")
)

//define gRPC error status codes
var ERRTYPE = map[error]int{
	DB_NAME_EMPTY:               3,
	DB_IDENTIFIER_EMPTY:         3,
	COLLECTION_NAME_EMPTY:       3,
	COLLECTION_IDENTIFIER_EMPTY: 3,
	NAMESPACE_IDENTIFIER_EMPTY:  3,
	KEY_EMPTY:                   3,
	EMPTY_KEY_CANNOT_BE_DELETED: 3,
	START_KEY_UNKNOWN:           3,
	START_OR_END_KEY_EMPTY:      3,
	IDENTIFIER_NOT_FOUND:        5,
}
