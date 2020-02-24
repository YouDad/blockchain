package core

import (
	"bytes"

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

func BytesToTxnOutput(bytes []byte) []TxnOutput {
	txnOutputs := []TxnOutput{}
	log.Err(utils.GetDecoder(bytes).Decode(&txnOutputs))

	return txnOutputs
}
