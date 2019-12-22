package core

import (
	"log"

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
	lastBlock := DeserializeBlock(bc.Blocks().GetLastest())
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

func (bc *Blockchain) AddBlock(data app.App) *Block {
	lastestBlock := DeserializeBlock(bc.Blocks().GetLastest())
	newBlock := NewBlock(data, lastestBlock.Hash)
	bc.Blocks().SetLastest(newBlock.Hash, newBlock.Serialize())
	return newBlock
}

func NewBlockchain() *Blockchain {
	if !utils.IsDatabaseExists() {
		log.Panicln("No existing blockchain found. Create one to continue.")
	}

	return &Blockchain{utils.OpenDatabase()}
}

func CreateBlockchain() *Blockchain {
	if utils.IsDatabaseExists() {
		log.Panicln("Blockchain existed, Create failed.")
	}

	db := utils.OpenDatabase()
	db.Blocks().Clear()
	genesis := NewBlock(coreConfig.GetGenesis(), make([]byte, 32))
	db.SetLastest(genesis.Hash, genesis.Serialize())
	return &Blockchain{db}
}
