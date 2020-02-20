package core

import (
	"errors"
	"log"

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
	iter.Blockchain.Blocks()
	blockBytes := iter.Get(iter.next)
	if len(blockBytes) == 0 {
		return nil
	}
	nextBlock = DeserializeBlock(blockBytes)
	iter.next = nextBlock.PrevBlockHash
	return nextBlock
}

func (bc *Blockchain) MineBlock(data CoinApp) *Block {
	lastestBlock := DeserializeBlock(bc.GetLastest())
	newBlock := NewBlock(data, lastestBlock.Hash, lastestBlock.Height+1)
	bc.SetLastest(newBlock.Hash, newBlock.Serialize())
	bc.SetByInt(newBlock.Height, newBlock.Serialize())
	return newBlock
}

func (bc *Blockchain) AddBlock(block *Block) {
	if block == nil {
		return
	}
	if bc.Blocks().Get(block.Hash) != nil {
		return
	}

	lastestBlock := DeserializeBlock(bc.GetLastest())
	if lastestBlock.Height < block.Height {
		bc.SetLastest(block.Hash, block.Serialize())
		bc.SetByInt(block.Height, block.Serialize())
	}
}

func IsBlockchainExists() bool {
	return utils.IsDatabaseExists(CoreConfig.DatabaseFile)
}

func NewBlockchai() *Blockchain {
	if !utils.IsDatabaseExists(CoreConfig.DatabaseFile) {
		log.Panicln("No existing blockchain found. Create one to continue.")
	}

	return &Blockchain{utils.OpenDatabase(CoreConfig.DatabaseFile)}
}

func CreateBlockchainFromGenesis(genesis *Block) *Blockchain {
	if utils.IsDatabaseExists(CoreConfig.DatabaseFile) {
		log.Panicln("Blockchain existed, Create failed.")
	}

	db := utils.OpenDatabase(CoreConfig.DatabaseFile)
	db.Blocks().Clear()
	db.SetGenesis(genesis.Hash, genesis.Serialize())
	db.SetByInt(genesis.Height, genesis.Serialize())
	return &Blockchain{db}
}

func CreateBlockchai() *Blockchain {
	if utils.IsDatabaseExists(CoreConfig.DatabaseFile) {
		log.Panicln("Blockchain existed, Create failed.")
	}

	db := utils.OpenDatabase(CoreConfig.DatabaseFile)
	db.Blocks().Clear()
	genesis := NewBlock(CoreConfig.GetGenesis(), make([]byte, 32), 1)
	db.SetGenesis(genesis.Hash, genesis.Serialize())
	db.SetByInt(genesis.Height, genesis.Serialize())
	return &Blockchain{db}
}

func (bc *Blockchain) GetBestHeight() int {
	return DeserializeBlock(bc.GetLastest()).Height
}

func (bc *Blockchain) GetBlock(hash []byte) (*Block, error) {
	value := bc.Blocks().Get(hash)
	if len(value) == 0 {
		return nil, errors.New("Block is not found.")
	}
	return DeserializeBlock(value), nil
}

func (bc *Blockchain) GetBlockHashes() (hashes [][]byte) {
	iter := bc.Begin()
	for {
		block := iter.Next()
		if block == nil {
			break
		}
		hashes = append(hashes, block.Hash)
	}
	return hashes
}
