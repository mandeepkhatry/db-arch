package def

/*
Package def defines constants, error messages and their status codes
*/

import (
	"errors"

	"google.golang.org/grpc/codes"
)

//ERROR MESSAGES
var (
	DB_NAME_EMPTY               error = errors.New("dbname empty")
	DB_DOES_NOT_EXIST           error = errors.New("database doesn't exist")
	DB_IDENTIFIER_EMPTY         error = errors.New("dbidentifier empty")
	COLLECTION_NAME_EMPTY       error = errors.New("collection name empty")
	COLLECTION_IDENTIFIER_EMPTY error = errors.New("collection identifier empty")
	COLLECTION_DOES_NOT_EXIST   error = errors.New("collection doesn't exist")
	NAMESPACE_IDENTIFIER_EMPTY  error = errors.New("namespace identifier empty")
	NAMES_CANNOT_BE_EMPTY       error = errors.New("database/collection/namespace names can't be empty")
	NAMESPACE_DOES_NOT_EXIST    error = errors.New("namespace doesn't exist")
	KEY_EMPTY                   error = errors.New("key is empty")
	EMPTY_KEY_CANNOT_BE_DELETED error = errors.New("can't delete empty key")
	START_OR_END_KEY_EMPTY      error = errors.New("start or end key is empty")
	START_KEY_UNKNOWN           error = errors.New("Can't scan from last without knowing startKey")
	IDENTIFIER_NOT_FOUND        error = errors.New("id not found for given db/collection/namespace")
	CONNECTION_NOT_ESTABLISHED  error = errors.New("connection to database not established")
)

//define gRPC error status codes
var ERRTYPE = map[error]codes.Code{
	DB_NAME_EMPTY:               codes.InvalidArgument,
	DB_IDENTIFIER_EMPTY:         codes.InvalidArgument,
	COLLECTION_NAME_EMPTY:       codes.InvalidArgument,
	COLLECTION_IDENTIFIER_EMPTY: codes.InvalidArgument,
	NAMESPACE_IDENTIFIER_EMPTY:  codes.InvalidArgument,
	NAMES_CANNOT_BE_EMPTY:       codes.InvalidArgument,
	KEY_EMPTY:                   codes.InvalidArgument,
	EMPTY_KEY_CANNOT_BE_DELETED: codes.InvalidArgument,
	START_KEY_UNKNOWN:           codes.InvalidArgument,
	START_OR_END_KEY_EMPTY:      codes.InvalidArgument,
	IDENTIFIER_NOT_FOUND:        codes.NotFound,
}
