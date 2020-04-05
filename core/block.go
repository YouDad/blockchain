package core

import (
	"encoding/json"

	"github.com/YouDad/blockchain/log"
	"github.com/YouDad/blockchain/types"
)

func BytesToBlock(bytes []byte) *types.Block {
	if bytes == nil {
		// log.SetCallerLevel(1)
		// log.Warnln("BytesToBlock parameter is nil")
		// log.SetCallerLevel(0)
		return nil
	}

	block := types.Block{}
	err := json.Unmarshal(bytes, &block)
	if err != nil {
		log.Warnf("%s\n", bytes)
		log.PrintStack()
	}
	log.Err(err)

	return &block
}
