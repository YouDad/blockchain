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
	log.NotImplement()
	return Block{}
}

func GetUTXOSet() *UTXOSet {
	return &UTXOSet{GetBlockchain()}
}

func (set *UTXOSet) Update(b *Block) {
	log.NotImplement()

}

func (set *UTXOSet) Reindex() {
	hashedUtxos := set.FindUTXO()
	set.SetTable("UTXOSet").Clear()

	for txnHash, utxos := range hashedUtxos {
		hash, err := hex.DecodeString(txnHash)
		log.Err(err)
		set.Set(hash, utils.Encode(utxos))
	}
}

func (set *UTXOSet) NewUTXOTransaction(from, to string, amount int) *Transaction {
	log.NotImplement()
	return nil
}

func (set *UTXOSet) FindUTXOByHash(pubKeyHash []byte) []TxnOutput {
	utxos := []TxnOutput{}

	set.SetTable("UTXOSet").Foreach(func(k, v []byte) bool {
		outs := BytesToTxnOutput(v)

		for _, out := range outs {
			if out.IsLockedWithKey(pubKeyHash) {
				utxos = append(utxos, out)
			}
		}
		return true
	})

	return utxos
}
