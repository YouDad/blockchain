package global

import (
	"sync"

	"github.com/YouDad/blockchain/log"
	"github.com/YouDad/blockchain/types"
)

type Mempool struct {
	pool map[int]map[[32]byte]types.Transaction
}

var instanceMempool Mempool
var onceMempool = sync.Once{}

func GetMempool() *Mempool {
	onceMempool.Do(func() {
		instanceMempool = Mempool{make(map[int]map[[32]byte]types.Transaction)}
	})
	return &instanceMempool
}

func (m Mempool) IsTxnExists(group int, txn types.Transaction) bool {
	_, ok := m.pool[group]
	if !ok {
		return false
	}

	var key [32]byte
	copy(key[0:32], txn.Hash())
	_, ok = m.pool[group][key]
	return ok
}

func (m *Mempool) AddTxnToMempool(group int, txn types.Transaction) {
	_, ok := m.pool[group]
	if !ok {
		m.pool[group] = make(map[[32]byte]types.Transaction)
	}

	log.Infof("AddTxnToMempool %s\n", txn.Hash())
	var key [32]byte
	copy(key[0:32], txn.Hash())
	m.pool[group][key] = txn
}

func (m Mempool) GetTxns(group int) []*types.Transaction {
	pool, ok := m.pool[group]
	if !ok {
		return nil
	}

	var ret []*types.Transaction
	for _, tx := range pool {
		ret = append(ret, &tx)
	}
	return ret
}

func (m Mempool) GetMempoolSize(group int) int {
	pool, ok := m.pool[group]
	if !ok {
		return 0
	}

	return len(pool)
}
