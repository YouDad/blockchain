package global

import (
	"sync"
)

type UTXOSetDB struct {
	IDatabase
}

var instanceUTXOSetDB *UTXOSetDB
var onceUTXOSetDB sync.Once

func GetUTXOSetDB() *UTXOSetDB {
	onceUTXOSetDB.Do(func() {
		instanceUTXOSetDB = &UTXOSetDB{getBoltDB("UTXOSet")}
	})
	return instanceUTXOSetDB
}
