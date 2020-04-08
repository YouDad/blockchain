package global

import (
	"sync"
)

type BlockheadsDB struct {
	IDatabase
}

var instanceBlockheadsDB *BlockheadsDB
var onceBlockheadsDB sync.Once

func GetBlockheadsDB() *BlockheadsDB {
	onceBlockheadsDB.Do(func() {
		instanceBlockheadsDB = &BlockheadsDB{getBoltDB("Blockheads")}
	})
	return instanceBlockheadsDB
}
