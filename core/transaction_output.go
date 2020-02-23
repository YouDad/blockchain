package core

import (
	"github.com/YouDad/blockchain/utils"
)

type TxnOutput struct {
	Value      int
	PubKeyHash []byte
}

func NewTxnOutput(address string, value int) *TxnOutput {
	pubKeyHash := utils.Base58Decode([]byte(address))
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-4]

	return &TxnOutput{value, pubKeyHash}
}
