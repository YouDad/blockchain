package global

import (
	"sync"

	"github.com/YouDad/blockchain/log"
	"github.com/YouDad/blockchain/types"
)

type Mempool struct {
	pool map[[32]byte]types.Transaction
}

var instanceMempool Mempool
var onceMempool = sync.Once{}

func GetMempool() Mempool {
	onceMempool.Do(func() {
		instanceMempool = Mempool{make(map[[32]byte]types.Transaction)}
	})
	return instanceMempool
}

func (m Mempool) IsTxnExists(txn types.Transaction) bool {
	var key [32]byte
	copy(key[0:32], txn.Hash())
	_, ok := m.pool[key]
	return ok
}

func (m *Mempool) AddTxnToMempool(txn types.Transaction) {
	log.Infof("AddTxnToMempool %x\n", txn.Hash)
	var key [32]byte
	copy(key[0:32], txn.Hash())
	m.pool[key] = txn
}

func (m Mempool) GetTxns() []*types.Transaction {
	var ret []*types.Transaction
	for _, tx := range m.pool {
		ret = append(ret, &tx)
	}
	return ret
}

func (m Mempool) GetMempoolSize() int {
	return len(m.pool)
}
