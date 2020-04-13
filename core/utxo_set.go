package core

import (
	"github.com/YouDad/blockchain/log"
)

type UTXOSet struct{}

func (set *UTXOSet) GetGenesis() Block {
	log.Errln("NotImplement")
	return Block{}
}

func NewUTXOSet() *UTXOSet {
	log.Errln("NotImplement")
	return nil
}

func (set *UTXOSet) Update(b *Block) {
	log.Errln("NotImplement")

}

func (set *UTXOSet) Reindex() {
	log.Errln("NotImplement")

}

func (set *UTXOSet) NewUTXOTransaction(from, to string, amount int) *Transaction {
	log.Errln("NotImplement")
	return nil
}
