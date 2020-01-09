package def

/*
Package def defines constants, error messages and their status codes
*/

const (
	META_DBIDENTIFIER                 = "meta:dbidentifier"
	META_COLLECTIONIDENTIFIER         = "meta:collectionidentifier"
	META_NAMESPACEIDENTIFIER          = "meta:namespaceidentifier"
	META_DBID                         = "meta:dbid:"
	META_COLLECTIONID                 = "meta:collectionid:"
	META_NAMESPACEID                  = "meta:namespaceid:"
	META_DB                           = "meta:db:"
	META_COLLECTION                   = "meta:collection:"
	META_NAMESPACE                    = "meta:namespace:"
	INDEX_KEY                         = "_index:"
	UNIQUE_ID                         = "_uniqueid:"
	UNIQUE_ID_INITIALCOUNT            = uint32(1)
	DBIDENTIFIER_INITIALCOUNT         = uint16(1)
	COLLECTIONIDENTIFIER_INITIALCOUNT = uint32(1)
	NAMESPACEIDENTIFIER_INITIALCOUNT  = uint32(1)
)
