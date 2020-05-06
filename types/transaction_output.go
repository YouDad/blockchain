package types

import (
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
	return out.PubKeyHash.Equal(pubKey.Hash())
}
