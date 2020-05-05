package core

import (
	"math/rand"
	"time"

	"github.com/YouDad/blockchain/log"
	"github.com/YouDad/blockchain/types"
	"github.com/YouDad/blockchain/utils"
)

func NewCoinbaseTxn(from string) *types.Transaction {
	randData := make([]byte, 32)
	rand.Seed(time.Now().UnixNano())
	for i := range randData {
		randData[i] = byte(rand.Int())
	}

	txn := types.Transaction{}

	txn.Vin = []types.TxnInput{{VoutIndex: -1, VoutValue: 50_000_000, PubKey: randData}}
	// Send $from 50BTC
	txn.Vout = []types.TxnOutput{*NewTxnOutput(from, 50_000_000)}
	return &txn
}

func BytesToTransaction(bytes []byte) *types.Transaction {
	txn := types.Transaction{}

	err := utils.Decode(bytes, &txn)
	if err != nil {
		log.Warn(err)
		log.Warnf("len=%d,bytes=%x", len(bytes), bytes)
		log.PrintStack()
	}
	return &txn
}
