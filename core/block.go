package core

import (
	"fmt"
	"time"

	"github.com/YouDad/blockchain/log"
	"github.com/YouDad/blockchain/types"
	"github.com/YouDad/blockchain/utils"
)

type BlockHead struct {
	Timestamp  int64
	PrevHash   types.HashValue
	Nonce      int64
	Height     int
	MerkleRoot types.HashValue
}

func (bh *BlockHead) Hash() types.HashValue {
	return utils.SHA256(*bh)
}

type Block struct {
	BlockHead
	Txns []*Transaction
}

func (b *Block) String() string {
	ret := fmt.Sprintf("Height: %d\n", b.Height)
	ret += fmt.Sprintf("Prev: %x\n", b.PrevHash)
	ret += fmt.Sprintf("Hash: %x\n", b.Hash())
	ret += fmt.Sprintf("Txns:\n")
	for i, txn := range b.Txns {
		ret += fmt.Sprintf("\tTxns[%d]:\n%s", i, txn.String())
	}

	return ret
}

func NewBlock(prev types.HashValue, height int, txns []*Transaction) *Block {
	block := &Block{
		BlockHead: BlockHead{
			Timestamp: time.Now().Unix(),
			PrevHash:  prev,
			Height:    height,
		},
		Txns: txns,
	}

	pow := NewPOW(block)
	nonce, _ := pow.Run()
	block.Nonce = nonce
	return block
}

func BytesToBlock(bytes []byte) *Block {
	if bytes == nil {
		log.SetCallerLevel(1)
		log.Warnln("BytesToBlock parameter is nil")
		log.SetCallerLevel(0)
		return nil
	}

	block := Block{}
	log.Err(utils.GetDecoder(bytes).Decode(&block))

	return &block
}
