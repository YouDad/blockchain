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

type TxnsDB struct {
	IDatabase
}

var instanceTxnsDB *TxnsDB
var onceTxnsDB sync.Once

func GetTxnsDB() *TxnsDB {
	onceTxnsDB.Do(func() {
		instanceTxnsDB = &TxnsDB{getBoltDB("Txns")}
	})
	return instanceTxnsDB
}
