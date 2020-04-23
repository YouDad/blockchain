package types

import "fmt"

type TxnInput struct {
	VoutHash  HashValue // 引用的交易哈希
	VoutIndex int       // 引用的交易在区块的位置
	VoutValue int64     // 被引用时的余额
	Signature Signature // 引用的签名
	PubKey    PublicKey // 被引用的公钥
}

func (in TxnInput) String() (ret string) {
	ret = "\n"
	ret += fmt.Sprintf("        VoutHash:  %s\n", in.VoutHash)
	ret += fmt.Sprintf("        VoutIndex: %d\n", in.VoutIndex)
	ret += fmt.Sprintf("        VoutValue: %d\n", in.VoutValue)
	ret += fmt.Sprintf("        Signature: %s\n", in.Signature)
	ret += fmt.Sprintf("        PubKey:    %s\n", in.PubKey)
	return ret
}
