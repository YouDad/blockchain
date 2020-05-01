package core

import (
	"encoding/json"
	"sync"

	"github.com/YouDad/blockchain/log"
	"github.com/YouDad/blockchain/types"
)

var mutexToBlock sync.Mutex

func BytesToBlock(bytes []byte) *types.Block {
	if bytes == nil {
		return nil
	}

	block := types.Block{}
	mutexToBlock.Lock()
	var err error
	defer func() {
		if r := recover(); r != nil {
			log.Warnln(r)
			defer func() {
				if r := recover(); r != nil {
					log.Errln(r)
				}
			}()
			err = json.Unmarshal(bytes, &block)
		}
	}()
	err = json.Unmarshal(bytes, &block)
	mutexToBlock.Unlock()
	if err != nil {
		log.Warn(err)
		log.Warnf("len=%d,bytes=%x", len(bytes), bytes)
		log.PrintStack()
	}
	log.Err(err)

	return &block
}
