package types

type TxnInput struct {
	VoutHash   HashValue // 引用的交易哈希
	VoutIndex  int       // 引用的交易在区块的位置
	VoutValue  int64     // 被引用时的余额
	Signature  Signature // 引用的签名
	PubKeyHash []byte    // 被引用的公钥
}