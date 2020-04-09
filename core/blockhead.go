package core

import (
	"github.com/YouDad/blockchain/global"
	"github.com/YouDad/blockchain/types"
	"github.com/YouDad/blockchain/utils"
)

type Blockhead struct {
	db    *global.BlockheadsDB
	group int
}

func (bh *Blockhead) Clear() {
	bh.db.Clear(bh.group)
}

func (bh *Blockhead) Get(key interface{}) (value []byte) {
	return bh.db.Get(bh.group, key)
}

func (bh *Blockhead) Set(key interface{}, value []byte) {
	bh.db.Set(bh.group, key, value)
}

func (bh *Blockhead) Delete(key interface{}) {
	bh.db.Delete(bh.group, key)
}

func (bh *Blockhead) Foreach(fn func(k, v []byte) bool) {
	bh.db.Foreach(bh.group, fn)
}
func GetBlockhead(group int) *Blockhead {
	return &Blockhead{global.GetBlockheadsDB(), group}
}

func (bh *Blockhead) AddBlockhead(block *types.Block) {
	if block.Verify() {
		bh.Set(block.Height, utils.Encode(block))
	}
}

func (bh *Blockhead) GetBlockheadByHeight(height int32) *types.Block {
	return BytesToBlock(bh.Get(height))
}
