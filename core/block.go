package core

import (
	"bytes"
	"fmt"
	"time"

	"github.com/YouDad/blockchain/log"
	"github.com/YouDad/blockchain/types"
	"github.com/YouDad/blockchain/utils"
)

type Block struct {
	Timestamp  int64
	PrevHash   types.HashValue
	Hash       types.HashValue
	Nonce      int64
	Height     int32
	MerkleRoot types.HashValue
	Txns       []*Transaction
}

func (b Block) String() string {
	ret := fmt.Sprintf("Height: %d\n", b.Height)
	ret += fmt.Sprintf("Prev: %x\n", b.PrevHash)
	ret += fmt.Sprintf("Hash: %x\n", b.Hash)
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
	log.Infof("Mining %s\n", block.String())
	nonce, hash := pow.Run()
	a := utils.SHA256(bytes.Join(
		[][]byte{
			pow.block.PrevHash[:],
			pow.block.MerkleRoot,
			utils.IntToBytes(pow.block.Timestamp),
			utils.IntToBytes(targetBits),
			utils.IntToBytes(nonce),
		},
		[]byte{},
	))
	log.Infoln(a, hash)
	block.Nonce = nonce
	block.Hash = hash
	log.Infof("Mined %s\n", block.String())
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
