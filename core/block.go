package core

import (
	"fmt"
	"time"

	"github.com/YouDad/blockchain/log"
	"github.com/YouDad/blockchain/types"
	"github.com/YouDad/blockchain/utils"
)

type Block struct {
	Timestamp  int64
	PrevHash   types.HashValue
	Nonce      int64
	Height     int32
	MerkleRoot types.HashValue
	Txns       []*Transaction
}

func (b Block) Hash() types.HashValue {
	// TODO: hash diff
	hashString := fmt.Sprintf("%d%d%d", b.Timestamp, b.Nonce, b.Height)
	hashString += fmt.Sprintf("%x%x", b.PrevHash, b.MerkleRoot)
	return utils.SHA256(hashString)
}

func (b Block) String() string {
	ret := fmt.Sprintf("Height: %d\n", b.Height)
	ret += fmt.Sprintf("Prev: %x\n", b.PrevHash)
	ret += fmt.Sprintf("Hash: %x\n", b.Hash())
	ret += fmt.Sprintf("Txns:\n")
	for i, txn := range b.Txns {
		ret += fmt.Sprintf("\tTxns[%d]:\n%s", i, txn.String())
	}

	return ret
}

func NewBlock(prev types.HashValue, height int32, txns []*Transaction) *Block {
	block := &Block{
		Timestamp: time.Now().Unix(),
		PrevHash:  prev,
		Height:    height,
		Txns:      txns,
	}
	var txs [][]byte

	for _, tx := range block.Txns {
		txs = append(txs, utils.Encode(tx))
	}
	mTree := NewMerkleTree(txs)
	block.MerkleRoot = mTree.RootNode.Data

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
