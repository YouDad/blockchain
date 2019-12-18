package core

import (
	"github.com/YouDad/blockchain/app"
	"github.com/YouDad/blockchain/utils"
)

type Blockchain struct {
	*utils.Database
}

type BlockchainIterator struct {
	*Blockchain
	next []byte
}

func (bc *Blockchain) Begin() (iter *BlockchainIterator) {
	lastBlock := DeserializeBlock(bc.GetLastest())
	return &BlockchainIterator{bc, lastBlock.Hash}
}

func (iter *BlockchainIterator) Next() (nextBlock *Block) {
	blockBytes := iter.Get(iter.next)
	if len(blockBytes) == 0 {
		return nil
	}
	nextBlock = DeserializeBlock(blockBytes)
	iter.next = nextBlock.PrevBlockHash
	return nextBlock
}

func (bc *Blockchain) AddBlock(data app.App) {
	lastestBlock := DeserializeBlock(bc.GetLastest())
	newBlock := NewBlock(data, lastestBlock.Hash)
	bc.SetLastest(newBlock.Hash, newBlock.Serialize())
}

func NewBlockchain() *Blockchain {
	db, exists := utils.OpenDatabase()
	if !exists {
		genesis := NewBlock(coreConfig.GetGenesis(), make([]byte, 32))
		db.SetLastest(genesis.Hash, genesis.Serialize())
	}
	return &Blockchain{db}
}
