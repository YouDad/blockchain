package rpc

import (
	"encoding/hex"

	"github.com/YouDad/blockchain/core"
)

var (
	mempool = make(map[string]*core.Transaction)
)

func isTransactionExists(tx *core.Transaction) bool {
	_, ok := mempool[hex.EncodeToString(tx.ID)]
	return ok
}

func addTransactionToMempool(tx *core.Transaction) {
	mempool[hex.EncodeToString(tx.ID)] = tx
}

func getTransactions() []*core.Transaction {
	var ret []*core.Transaction
	for _, tx := range mempool {
		ret = append(ret, tx)
	}
	return ret
}

func getMempoolSize() int {
	return len(mempool)
}
