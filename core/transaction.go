package core

import (
	"github.com/YouDad/blockchain/log"
	"github.com/YouDad/blockchain/types"
)

type Transaction struct {
	Hash types.HashValue
	Vin  []TxnInput
	Vout []TxnOutput
}

func NewCoinbaseTxn(from string) *Transaction {
	log.Errln("NotImplement")
	return nil
}
