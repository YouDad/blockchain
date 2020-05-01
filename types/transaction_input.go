package types

import "github.com/YouDad/blockchain/utils"

type TxnInput struct {
	VoutHash  HashValue // 引用的交易哈希
	VoutIndex int       // 引用的交易在区块的位置
	VoutValue int64     // 被引用时的余额
	Signature Signature // 引用的签名
	PubKey    PublicKey // 被引用的公钥
}

func (in TxnInput) String() string {
	return string(utils.Encode(in))
}
