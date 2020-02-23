package core

import (
	"encoding/hex"

	"github.com/YouDad/blockchain/log"
	"github.com/YouDad/blockchain/utils"
)

type UTXOSet struct {
	*Blockchain
}

func (set *UTXOSet) GetGenesis() Block {
	log.Errln("NotImplement")
	return Block{}
}

func NewUTXOSet() *UTXOSet {
	return &UTXOSet{GetBlockchain()}
}

func (set *UTXOSet) Update(b *Block) {
	log.Errln("NotImplement")

}

func (set *UTXOSet) Reindex() {
	set.SetTable("UTXOSet").Clear()
	UTXO := set.FindUTXO()

	for txnHash, outs := range UTXO {
		hash, err := hex.DecodeString(txnHash)
		log.Err(err)
		set.Set(hash, utils.Encode(outs))
	}
}

func (set *UTXOSet) NewUTXOTransaction(from, to string, amount int) *Transaction {
	log.Errln("NotImplement")
	return nil
}
