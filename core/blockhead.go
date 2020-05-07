package core

import (
	"sync"

	"github.com/YouDad/blockchain/global"
	"github.com/YouDad/blockchain/types"
	"github.com/YouDad/blockchain/utils"
)

type Blockhead struct {
	db    *global.BlockheadsDB
	group int
}

func (bh *Blockhead) clear() {
	bh.db.Clear(bh.group)
}

func (bh *Blockhead) get(key interface{}) (value []byte) {
	return bh.db.Get(bh.group, key)
}

func (bh *Blockhead) set(key interface{}, value []byte) {
	bh.db.Set(bh.group, key, value)
}

func (bh *Blockhead) Delete(key interface{}) {
	bh.db.Delete(bh.group, key)
}

func (bh *Blockhead) foreach(fn func(k, v []byte) bool) {
	bh.db.Foreach(bh.group, fn)
}
func GetBlockhead(group int) *Blockhead {
	return &Blockhead{global.GetBlockheadsDB(), group}
}

func (bh *Blockhead) GetLastest() *types.Block {
	return BytesToBlock(bh.get("lastest"))
}

var (
	onceBlockheadGetHeight      = make(map[int]*sync.Once)
	cacheBlockheadHeight        = make(map[int]int32)
	mutexBlockheadHeight        sync.Mutex
	mutexOnceBlockheadGetHeight sync.Mutex
)

func (bh *Blockhead) GetHeight() int32 {
	mutexOnceBlockheadGetHeight.Lock()
	_, ok := onceBlockheadGetHeight[bh.group]
	if !ok {
		onceBlockheadGetHeight[bh.group] = &sync.Once{}
	}

	onceBlockheadGetHeight[bh.group].Do(func() {
		lastest := bh.GetLastest()
		mutexBlockheadHeight.Lock()
		if lastest == nil {
			cacheBlockheadHeight[bh.group] = -1
		} else {
			cacheBlockheadHeight[bh.group] = lastest.Height
		}
		mutexBlockheadHeight.Unlock()
	})
	mutexOnceBlockheadGetHeight.Unlock()
	mutexBlockheadHeight.Lock()
	defer mutexBlockheadHeight.Unlock()
	return cacheBlockheadHeight[bh.group]
}

func (bh *Blockhead) AddBlockhead(block *types.Block) bool {
	if block == nil || bh.get(block.Height) != nil ||
		!block.Verify() || bh.GetHeight()+1 != block.Height {
		return false
	}

	txns := block.Txns
	block.Txns = nil
	bytes := utils.Encode(block)

	mutexBlockheadHeight.Lock()
	cacheBlockheadHeight[bh.group] = block.Height
	mutexBlockheadHeight.Unlock()

	bh.set(block.Height, bytes)
	bh.set("lastest", bytes)

	block.Txns = txns
	return true
}

func (bh *Blockhead) GetBlockheadByHeight(height int32) *types.Block {
	return BytesToBlock(bh.get(height))
}
