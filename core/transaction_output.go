package core

import (
	"encoding/json"

	"github.com/YouDad/blockchain/log"
	"github.com/YouDad/blockchain/types"
	"github.com/YouDad/blockchain/utils"
)

func NewTxnOutput(address string, value int64) *types.TxnOutput {
	pubKeyHash := utils.Base58Decode([]byte(address))
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-4]

	return &types.TxnOutput{
		Value:      value,
		PubKeyHash: pubKeyHash,
	}
}

func BytesToTxnOutputs(bytes []byte) []types.TxnOutput {
	txnOutputs := []types.TxnOutput{}
	err := json.Unmarshal(bytes, &txnOutputs)
	if err != nil {
		log.Warnf("%s\n", bytes)
		log.PrintStack()
	}

	return txnOutputs
}
