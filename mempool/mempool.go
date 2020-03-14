package mempool

import (
	"github.com/YouDad/blockchain/core"
	"github.com/YouDad/blockchain/log"
)

var (
	mempool = make(map[[32]byte]core.Transaction)
)

func IsTxnExists(txn core.Transaction) bool {
	var key [32]byte
	copy(key[0:32], txn.Hash)
	_, ok := mempool[key]
	return ok
}

func AddTxnToMempool(txn core.Transaction) {
	log.Infof("AddTxnToMempool %x\n", txn.Hash)
	var key [32]byte
	copy(key[0:32], txn.Hash)
	mempool[key] = txn
}

func GetTxns() []*core.Transaction {
	var ret []*core.Transaction
	for _, tx := range mempool {
		ret = append(ret, &tx)
	}
	return ret
}

func GetMempoolSize() int {
	return len(mempool)
}
