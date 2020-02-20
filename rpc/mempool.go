package rpc

import (
	"encoding/hex"

	"github.com/YouDad/blockchain/coin_core"
)

var (
	mempool = make(map[string]*coin_core.Transaction)
)

func isTransactionExists(tx *coin_core.Transaction) bool {
	_, ok := mempool[hex.EncodeToString(tx.ID)]
	return ok
}

func addTransactionToMempool(tx *coin_core.Transaction) {
	mempool[hex.EncodeToString(tx.ID)] = tx
}

func getTransactions() []*coin_core.Transaction {
	var ret []*coin_core.Transaction
	for _, tx := range mempool {
		ret = append(ret, tx)
	}
	return ret
}

func getMempoolSize() int {
	return len(mempool)
}
