package kvstore

import (
	"db-arch/server/internal/kvstore/tikvdb"
	"db-arch/server/io"
)

//NewTiKVFactory returns tikv storeclient as io.Store
func NewTiKVFactory(pdAddr []string, dbDir string) io.Store {
	tikv := &tikvdb.StoreClient{}
	err := tikv.NewClient(pdAddr)
	if err != nil {
		panic(err)
	}
	return tikv
}

//NewBadgerFactory returns badgerdb storeclient as io.Store
func NewBadgerFactory(pdAddr []string, dbDir string) io.Store {
	badger := &badgerdb.StoreClient{}
	err := badger.NewClient(pdAddr, dbDir)
	if err != nil {
		panic(err)
	}
	return badger
}
