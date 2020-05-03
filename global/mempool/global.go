package mempool

import "github.com/YouDad/blockchain/types"

func AddTxn(group int, txn types.Transaction) {
	m := GetMempool(group)
	m.AddTxn(txn)
	m.Release()
}

func GetMempoolSize(group int) int {
	m := GetMempool(group)
	ret := m.GetMempoolSize()
	m.Release()
	return ret
}

func GetTxn(group int, hash types.HashValue) (*types.Transaction, error) {
	m := GetMempool(group)
	ret, err := m.GetTxn(hash)
	m.Release()
	return ret, err
}

func GetTxns(group int) []*types.Transaction {
	m := GetMempool(group)
	ret := m.GetTxns()
	m.Release()
	return ret
}

func ExpandTxnOutput(group int, out types.TxnOutput, hash types.HashValue, index int) (
	outs []*types.TxnOutput, hashs []types.HashValue, indexs []int) {
	m := GetMempool(group)
	outs, hashs, indexs = m.ExpandTxnOutput(out, hash, index)
	m.Release()
	return outs, hashs, indexs
}
