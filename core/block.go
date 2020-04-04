package core

import (
	"encoding/json"
	"time"

	"github.com/YouDad/blockchain/log"
	"github.com/YouDad/blockchain/types"
	"github.com/YouDad/blockchain/utils"
)

func NewBlock(group int, prev types.HashValue, target float64, height int32, txns []*types.Transaction) *types.Block {
	block := &types.Block{
		BlockHeader: types.BlockHeader{
			Timestamp: time.Now().UnixNano(),
			PrevHash:  prev,
			Target:    target,
			Height:    height,
		},
		Txns: txns,
	}
	var txs [][]byte

	for _, tx := range block.Txns {
		txs = append(txs, utils.Encode(tx))
	}
	mTree := NewMerkleTree(txs)
	block.MerkleRoot = mTree.RootNode.Data

	pow := NewPOW(block)
	nonce, hash := pow.Run(group)
	if hash == nil {
		return nil
	}
	block.Nonce = nonce
	log.Debugf("Mined %s\n", block.String())
	return block
}

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
