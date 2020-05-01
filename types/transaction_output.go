package types

import (
	"bytes"

	"github.com/YouDad/blockchain/utils"
)

type TxnOutput struct {
	Value      int64
	PubKeyHash HashValue
}

func (out TxnOutput) String() string {
	return string(utils.Encode(out))
}

func (out *TxnOutput) IsLockedWithKey(pubKey PublicKey) bool {
	return bytes.Compare(out.PubKeyHash, pubKey.Hash()) == 0
}
