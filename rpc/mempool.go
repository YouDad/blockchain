package rpc

import (
	"encoding/hex"

	coin_core "github.com/YouDad/blockchain/app/coin/core"
)

var (
	mempool = make(map[string]*coin_core.Transaction)
)

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
