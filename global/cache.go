package global

import "github.com/YouDad/blockchain/types"

var groupCache = make(map[int]map[interface{}]*types.Block)

func SetBlock(group int, key interface{}, block *types.Block) {
	blockCache, ok := groupCache[group]
	if !ok {
		groupCache[group] = make(map[interface{}]*types.Block)
		blockCache = groupCache[group]
	}
	blockCache[key] = block
}

func GetBlock(group int, key interface{}) (*types.Block, bool) {
	blockCache, ok := groupCache[group]
	if !ok {
		return nil, false
	}
	block, ok := blockCache[key]
	if !ok {
		return nil, false
	}
	return block, ok
}
