package global

import (
	"errors"
	"fmt"
	"sync"

	"github.com/YouDad/blockchain/log"
	"github.com/YouDad/blockchain/types"
)

type Mempool map[[32]byte]types.Transaction

type Mempools map[int]Mempool

var instanceMempool Mempools
var onceMempool = sync.Once{}

func GetMempool(group int) Mempool {
	onceMempool.Do(func() {
		instanceMempool = make(Mempools)
	})

	_, ok := instanceMempool[group]
	if !ok {
		instanceMempool[group] = make(Mempool)
	}
	return instanceMempool[group]
}

func (m Mempool) IsTxnExists(txn types.Transaction) bool {
	var key [32]byte
	copy(key[0:32], txn.Hash())
	_, ok := m[key]
	return ok
}

func (m Mempool) AddTxnToMempool(txn types.Transaction) {
	log.Infof("AddTxnToMempool %s\n", txn.Hash())
	var key [32]byte
	copy(key[0:32], txn.Hash())
	m[key] = txn
}

func (m Mempool) GetTxns() []*types.Transaction {
	var ret []*types.Transaction
	i := 0
	for _, txn := range m {
		ret = append(ret, &txn)
		i++
		if i == 2000 {
			break
		}
	}
	return ret
}

func (m Mempool) GetMempoolSize() int {
	return len(m)
}

func (m Mempool) FindTxn(hash types.HashValue) (*types.Transaction, error) {
	for _, txn := range m {
		if txn.Hash().Equal(hash) {
			return &txn, nil
		}
	}
	return nil, errors.New(fmt.Sprintf("Transaction is not found, %s", hash))
}

func (m Mempool) FindTxnOutput(out types.TxnOutput, hash types.HashValue, index int) (
	outs []*types.TxnOutput, hashs []types.HashValue, indexs []int, err error) {
	outs = append(outs, &out)
	hashs = append(hashs, hash)
	indexs = append(indexs, index)

	target := 0
	for target != len(outs) {
		outs = outs[target:]
		hashs = hashs[target:]
		indexs = indexs[target:]
		target = len(outs)

		for i := 0; i < target; i++ {
			for _, txn := range m {
				for _, in := range txn.Vin {
					if !(outs[i].IsLockedWithKey(in.PubKeyHash) &&
						outs[i].Value == in.VoutValue && hashs[i].Equal(in.VoutHash)) {
						continue
					}

					for index, out := range txn.Vout {
						if !out.IsLockedWithKey(in.PubKeyHash) {
							continue
						}

						outs = append(outs, &out)
						hashs = append(hashs, txn.Hash())
						indexs = append(indexs, index)
					}
				}
			}
		}
	}
	return outs, hashs, indexs, err
}
