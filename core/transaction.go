package core

import (
	"crypto/rand"
	"encoding/json"
	"fmt"

	"github.com/YouDad/blockchain/log"
	"github.com/YouDad/blockchain/types"
	"github.com/YouDad/blockchain/utils"
)

func NewCoinbaseTxn(from string) *types.Transaction {
	randData := make([]byte, 20)
	_, err := rand.Read(randData)
	log.Err(err)
	data := fmt.Sprintf("%x", randData)

	txn := types.Transaction{}

	txn.Vin = []types.TxnInput{types.TxnInput{VoutIndex: -1, PubKeyHash: []byte(data)}}
	// Send $from 50BTC
	txn.Vout = []types.TxnOutput{*NewTxnOutput(from, 50_000_000)}

	txn.Hash = utils.SHA256(&txn)
	return &txn
}

func BytesToTransaction(bytes []byte) *types.Transaction {
	txn := types.Transaction{}
	err := json.Unmarshal(bytes, &txn)
	if err != nil {
		log.Warnf("%s\n", bytes)
		log.PrintStack()
	}
	return &txn
}
