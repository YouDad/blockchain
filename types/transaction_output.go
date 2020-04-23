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
	ret += fmt.Sprintf("\t\tValue: %d\n", out.Value)
	ret += fmt.Sprintf("\t\tPubKeyHash: %s\n", out.PubKeyHash)
	return ret
}

func (out *TxnOutput) IsLockedWithKey(pubKeyHash []byte) bool {
	return bytes.Compare(out.PubKeyHash, pubKeyHash) == 0
}
