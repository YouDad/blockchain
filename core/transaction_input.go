package core

import "github.com/YouDad/blockchain/types"

type TxnInput struct {
	VoutHash   types.HashValue
	VoutIndex  int
	Signature  types.Signature
	PubKeyHash []byte
}
