package types

import "fmt"

type TxnInput struct {
	VoutHash   HashValue // 引用的交易哈希
	VoutIndex  int       // 引用的交易在区块的位置
	VoutValue  int64     // 被引用时的余额
	Signature  Signature // 引用的签名
	PubKeyHash HashValue // 被引用的公钥
}

func (in TxnInput) String() (ret string) {
	ret = "\n"
	ret += fmt.Sprintf("\t\tVoutHash:   %s\n", in.VoutHash)
	ret += fmt.Sprintf("\t\tVoutIndex:  %d\n", in.VoutIndex)
	ret += fmt.Sprintf("\t\tVoutValue:  %d\n", in.VoutValue)
	ret += fmt.Sprintf("\t\tSignature:  %s\n", in.Signature)
	ret += fmt.Sprintf("\t\tPubKeyHash: %s\n", in.PubKeyHash)
	return ret
}
