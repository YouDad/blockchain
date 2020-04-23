package types

import (
	"bytes"
	"fmt"
)

type TxnOutput struct {
	Value      int64
	PubKeyHash HashValue
}

func (out TxnOutput) String() (ret string) {
	ret = "\n"
	ret += fmt.Sprintf("        Value: %d\n", out.Value)
	ret += fmt.Sprintf("        PubKeyHash: %s\n", out.PubKeyHash)
	return ret
}

func (out *TxnOutput) IsLockedWithKey(pubKey PublicKey) bool {
	return bytes.Compare(out.PubKeyHash, pubKey.Hash()) == 0
}
