package types

import "bytes"

type TxnOutput struct {
	Value      int64
	PubKeyHash []byte
}

func (out *TxnOutput) IsLockedWithKey(pubKeyHash []byte) bool {
	return bytes.Compare(out.PubKeyHash, pubKeyHash) == 0
}
