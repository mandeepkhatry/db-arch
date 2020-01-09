package kvstore

import (
	"db-arch/server/internal/kvstore/tikvdb"
	"db-arch/server/io"
)

//NewTiKVFactory returns tikv storeclient as io.Store
func NewTiKVFactory(pdAddr []string) io.Store {
	tikv := &tikvdb.StoreClient{}
	err := tikv.NewClient(pdAddr)
	if err != nil {
		panic(err)
	}
	return tikv
}
