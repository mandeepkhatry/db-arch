package kvstore

//IndexDocument indexes document in batch
func (s *StoreClient) IndexDocument(dbID []byte,collectionID []byte,
	namespaceID []byte,uniqueID []byte,
	data map[string][]byte,indices []string)error{

	/*

	1. Find type of data
	2. Generate index key
	3. Read db
	4. Append unique_id using roaring bitmap
	5. Write to db in batch

	 */


}
