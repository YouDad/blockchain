package global

import (
	"sync"
)

type BlocksDB struct {
	IDatabase
}

var instanceBlocksDB *BlocksDB
var onceBlocksDB sync.Once

func GetBlocksDB() *BlocksDB {
	onceBlocksDB.Do(func() {
		instanceBlocksDB = &BlocksDB{getBoltDB("Blocks")}
	})
	return instanceBlocksDB
}
