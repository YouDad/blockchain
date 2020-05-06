package mempool

import (
	"errors"
	"fmt"
	"sync"

	"github.com/YouDad/blockchain/global"
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
	group = group % global.MaxGroupNum
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

func (m Mempool) Release() {
	instanceMempool.mutex[m.group].Unlock()
}

func (m Mempool) AddTxn(txn types.Transaction) {
	m.m[txn.Hash().Key()] = txn
}

func (m Mempool) GetTxn(hash types.HashValue) (*types.Transaction, error) {
	txn, ok := m.m[hash.Key()]
	if ok {
		return &txn, nil
	}
	return nil, errors.New(fmt.Sprintf("Transaction is not found, %s", hash))
}

func (m Mempool) Delete(hash types.HashValue) {
	delete(m.m, hash.Key())
	log.SetCallerLevel(1)
	log.Debugln("Mempool Txn Delete", hash)
	log.SetCallerLevel(0)
}

func (m Mempool) GetTxns() []*types.Transaction {
	// 拓扑排序结构和函数定义
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

	// 将目标映射到整数域
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

	// 构建拓扑排序图
	for _, txn := range m.m {
		for _, vin := range txn.Vin {
			prevTxn, ok := m.m[vin.VoutHash.Key()]
			if ok {
				j := keyToId[prevTxn.Hash().Key()]
				addedge(j, keyToId[txn.Hash().Key()])
			}
		}
	}

	// 选出入度为0的节点作为初始节点
	var nodes []int
	lenBefore := len(nodes)
	for node, cnt := range indeg {
		if cnt == 0 {
			nodes = append(nodes, node)
		}
	}
	lenAfter := len(nodes)

	// 开始排序
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

	// 收集返回值及删除有效节点
	var ret []*types.Transaction
	i = 0
	for _, node := range nodes {
		if i < 50 {
			ret = append(ret, idToTxn[node])
		}
		delete(idToTxn, node)
		i++
	}

	// 删除拓扑排序中的环结构
	for _, txn := range idToTxn {
		delete(m.m, txn.Hash().Key())
	}

	return ret
}

func (m Mempool) GetMempoolSize() int {
	return len(m.m)
}

func (m Mempool) ExpandTxnOutput(out types.TxnOutput, hash types.HashValue, index int) (
	outs []*types.TxnOutput, hashs []types.HashValue, indexs []int) {
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
