package global

import (
	"errors"
	"fmt"
	"sync"

	"github.com/YouDad/blockchain/log"
	"github.com/YouDad/blockchain/types"
)

type Mempool struct {
	m     map[[32]byte]types.Transaction
	group int
}
type Mempools struct {
	mempool map[int]Mempool
	mutex   map[int]*sync.Mutex
}

var instanceMempool Mempools
var onceMempool = sync.Once{}

func GetMempool(group int) Mempool {
	onceMempool.Do(func() {
		instanceMempool.mempool = make(map[int]Mempool)
		instanceMempool.mutex = make(map[int]*sync.Mutex)
	})

	_, ok := instanceMempool.mempool[group]
	if !ok {
		instanceMempool.mempool[group] = Mempool{make(map[[32]byte]types.Transaction), group}
		instanceMempool.mutex[group] = &sync.Mutex{}
	}
	instanceMempool.mutex[group].Lock()
	return instanceMempool.mempool[group]
}

func (m Mempool) release() {
	instanceMempool.mutex[m.group].Unlock()
}

func (m Mempool) AddTxn(txn types.Transaction) {
	defer m.release()
	log.Infof("AddTxn %s\n", txn.Hash())
	var key [32]byte
	copy(key[0:32], txn.Hash())
	m.m[key] = txn
}

func (m Mempool) GetTxn(hash types.HashValue) (*types.Transaction, error) {
	defer m.release()
	txn, ok := m.m[hash.Key()]
	if ok {
		return &txn, nil
	}
	return nil, errors.New(fmt.Sprintf("Transaction is not found, %s", hash))
}

func (m Mempool) Delete(hash types.HashValue) {
	defer m.release()
	delete(m.m, hash.Key())
}

func (m Mempool) GetTxns() []*types.Transaction {
	defer m.release()

	// 拓扑排序
	type edge struct {
		dest int
		next *edge
	}
	var e []edge
	head := make(map[int]*edge)
	indeg := make(map[int]int)
	addedge := func(u, v int) {
		e = append(e, edge{v, head[u]})
		head[u] = &e[len(e)-1]
		_, ok := indeg[u]
		if !ok {
			indeg[u] = 0
		}
		indeg[v] += 1
	}

	keyToId := make(map[[32]byte]int)
	idToTxn := make(map[int]*types.Transaction)
	i := 0
	for _, txn := range m.m {
		keyToId[txn.Hash().Key()] = i
		copyTxn := txn
		idToTxn[i] = &copyTxn
		indeg[i] = 0
		i++
	}

	for _, txn := range m.m {
		for _, vin := range txn.Vin {
			prevTxn, ok := m.m[vin.VoutHash.Key()]
			if ok {
				j := keyToId[prevTxn.Hash().Key()]
				addedge(j, keyToId[txn.Hash().Key()])
			}
		}
	}

	var nodes []int
	lenBefore := len(nodes)
	for node, cnt := range indeg {
		if cnt == 0 {
			nodes = append(nodes, node)
		}
	}
	lenAfter := len(nodes)

	for lenAfter != lenBefore {
		for i := lenBefore; i < lenAfter; i++ {
			u := nodes[i]
			delete(indeg, u)
			for ptr := head[u]; ptr != nil; ptr = ptr.next {
				indeg[ptr.dest] -= 1
				if indeg[ptr.dest] == 0 {
					nodes = append(nodes, ptr.dest)
				}
			}
		}
		lenBefore = lenAfter
		lenAfter = len(nodes)
	}

	var ret []*types.Transaction
	i = 0
	for _, node := range nodes {
		ret = append(ret, idToTxn[node])
		i++
		if i == 50 {
			break
		}
	}
	return ret
}

func (m Mempool) GetMempoolSize() int {
	defer m.release()
	return len(m.m)
}

func (m Mempool) ExpandTxnOutput(out types.TxnOutput, hash types.HashValue, index int) (
	outs []*types.TxnOutput, hashs []types.HashValue, indexs []int) {
	defer m.release()
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
			for _, txn := range m.m {
				for _, in := range txn.Vin {
					if !(outs[i].IsLockedWithKey(in.PubKey) &&
						outs[i].Value == in.VoutValue && hashs[i].Equal(in.VoutHash)) {
						continue
					}

					for index, out := range txn.Vout {
						if !out.IsLockedWithKey(in.PubKey) {
							continue
						}

						copyOut := out
						outs = append(outs, &copyOut)
						hashs = append(hashs, txn.Hash())
						indexs = append(indexs, index)
					}
				}
			}
		}
	}
	return outs, hashs, indexs
}
