package core

import (
	"sync"

	"github.com/YouDad/blockchain/log"
	"github.com/YouDad/blockchain/types"
	"github.com/YouDad/blockchain/utils"
)

var mutexToBlock sync.Mutex

func BytesToBlock(bytes []byte) *types.Block {
	if bytes == nil {
		return nil
	}

	block := types.Block{}
	err := utils.Decode(bytes, &block)
	if err != nil {
		log.Warn(err)
		log.Warnf("len=%d,bytes=%x", len(bytes), bytes)
		log.PrintStack()
	}
	log.Err(err)

	return &block
}
