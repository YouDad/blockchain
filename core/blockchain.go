package core

import (
	"github.com/YouDad/blockchain/log"
	"github.com/YouDad/blockchain/storage"
)

type Blockchain struct {
	storage.Database
}

func CreateBlockchain(minerAddress string) {
	log.Errln("NotImplement")
}

func CreateBlockchainFromGenesis(b *Block) *Blockchain {
	log.Errln("NotImplement")
	return nil
}

func IsExists() bool {
	log.Errln("NotImplement")
	return true
}

func GetBlockchain() *Blockchain {
	log.Errln("NotImplement")
	return &Blockchain{storage.GetDatabase()}
}

func (bc *Blockchain) GetGenesis() *Block {
	log.Errln("NotImplement")
	return nil
}

func (bc *Blockchain) GetHeight() int {
	log.Errln("NotImplement")
	return 0
}

func (bc *Blockchain) AddBlock(b *Block) {
	log.Errln("NotImplement")

}

func (bc *Blockchain) MineBlock(txns []*Transaction) *Block {
	log.Errln("NotImplement")
	return nil
}
