package core

import (
	"encoding/json"

	"github.com/YouDad/blockchain/log"
	"github.com/YouDad/blockchain/types"
)

func BytesToBlock(bytes []byte) *types.Block {
	if bytes == nil {
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
