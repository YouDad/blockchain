package core

import (
	"crypto/rand"
	"encoding/json"
	"fmt"

	"github.com/YouDad/blockchain/log"
	"github.com/YouDad/blockchain/types"
)

func NewCoinbaseTxn(from string) *types.Transaction {
	randData := make([]byte, 20)
	_, err := rand.Read(randData)
	log.Err(err)
	data := fmt.Sprintf("%x", randData)

	txn := types.Transaction{}

	txn.Vin = []types.TxnInput{{VoutIndex: -1, PubKeyHash: []byte(data)}}
	// Send $from 50BTC
	txn.Vout = []types.TxnOutput{*NewTxnOutput(from, 50_000_000)}
	return &txn
}

func BytesToTransaction(bytes []byte) *types.Transaction {
	txn := types.Transaction{}
	err := json.Unmarshal(bytes, &txn)
	if err != nil {
		log.Warn(err)
		log.Warnf("len=%d,bytes=%x", len(bytes), bytes)
		log.PrintStack()
	}
	return &txn
}
