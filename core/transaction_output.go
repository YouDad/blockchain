package core

import (
	"bytes"
	"encoding/json"

	"github.com/YouDad/blockchain/log"
	"github.com/YouDad/blockchain/utils"
)

type TxnOutput struct {
	Value      int64
	PubKeyHash []byte
}

func NewTxnOutput(address string, value int64) *TxnOutput {
	pubKeyHash := utils.Base58Decode([]byte(address))
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-4]

	return &TxnOutput{value, pubKeyHash}
}

func (out *TxnOutput) IsLockedWithKey(pubKeyHash []byte) bool {
	return bytes.Compare(out.PubKeyHash, pubKeyHash) == 0
}

func BytesToTxnOutputs(bytes []byte) []TxnOutput {
	txnOutputs := []TxnOutput{}
	err := json.Unmarshal(bytes, &txnOutputs)
	if err != nil {
		log.Tracef("%s\n", bytes)
		log.PrintStack()
	}

	return txnOutputs
}
