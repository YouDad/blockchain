package core

import (
	"encoding/hex"

	"github.com/YouDad/blockchain/conf"
	"github.com/YouDad/blockchain/log"
	"github.com/YouDad/blockchain/types"
	"github.com/YouDad/blockchain/utils"
	"github.com/YouDad/blockchain/wallet"
)

type UTXOSet struct {
	*Blockchain
}

func GetUTXOSet() *UTXOSet {
	return &UTXOSet{GetBlockchain()}
}

func (set *UTXOSet) Update(b *Block) {
	for _, tx := range b.Txns {
		if tx.IsCoinbase() == false {
			for _, vin := range tx.Vin {
				updatedOuts := []TxnOutput{}
				outsBytes := set.SetTable(conf.UTXOSET).Get(vin.VoutHash)
				outs := BytesToTxnOutputs(outsBytes)

				for outIdx, out := range outs {
					if outIdx != vin.VoutIndex {
						updatedOuts = append(updatedOuts, out)
					}
				}

				if len(updatedOuts) == 0 {
					set.SetTable(conf.UTXOSET).Delete(vin.VoutHash)
				} else {
					set.SetTable(conf.UTXOSET).Set(vin.VoutHash, utils.Encode(updatedOuts))
				}
			}
		}

		newOutputs := []TxnOutput{}
		for _, out := range tx.Vout {
			newOutputs = append(newOutputs, out)
		}
		set.SetTable(conf.UTXOSET).Set(tx.Hash, utils.Encode(newOutputs))
	}
}

func (set *UTXOSet) Reindex() {
	hashedUtxos := set.FindUTXO()
	set.SetTable(conf.UTXOSET).Clear()

	for txnHash, utxos := range hashedUtxos {
		hash, err := hex.DecodeString(txnHash)
		log.Err(err)
		set.Set(hash, utils.Encode(utxos))
	}
}

func (set *UTXOSet) NewUTXOTransaction(from, to string, amount int64) *Transaction {
	var ins []TxnInput
	var outs []TxnOutput

	wallets, err := wallet.NewWallets()
	log.Err(err)

	srcWallet, have := wallets.GetWallet(from)
	if !have {
		log.Errf("You haven't %s's PrivateKey", from)
	}
	pubKeyHash := wallet.HashPubKey(srcWallet.PublicKey)
	acc, utxos := set.FindUTXOs(pubKeyHash, amount)

	if acc < amount {
		log.Errln("Not enough BTC")
	}

	for txnHash, outIdxs := range utxos {
		txnId, err := hex.DecodeString(txnHash)
		log.Err(err)

		for _, outIdx := range outIdxs {
			ins = append(ins, TxnInput{
				txnId, outIdx, nil, srcWallet.PublicKey})
		}
	}

	outs = append(outs, *NewTxnOutput(to, amount))
	if acc > amount {
		outs = append(outs, *NewTxnOutput(from, acc-amount))
	}

	txn := Transaction{nil, ins, outs}
	txn.Hash = utils.SHA256(&txn)
	set.SignTransaction(&txn, srcWallet.PrivateKey)
	return &txn
}

func (set *UTXOSet) FindUTXOByHash(pubKeyHash []byte) []TxnOutput {
	utxos := []TxnOutput{}

	set.SetTable(conf.UTXOSET).Foreach(func(k, v []byte) bool {
		outs := BytesToTxnOutputs(v)

		for _, out := range outs {
			if out.IsLockedWithKey(pubKeyHash) {
				utxos = append(utxos, out)
			}
		}
		return true
	})

	return utxos
}

func (set *UTXOSet) FindUTXOs(pubKeyHash types.HashValue, amount int64) (int64, map[string][]int) {
	hashedUTxnoutIdxs := make(map[string][]int)
	var accumulated int64 = 0

	set.SetTable(conf.UTXOSET).Foreach(func(k, v []byte) bool {
		txnHash := hex.EncodeToString(k)
		outs := BytesToTxnOutputs(v)

		for outIdx, out := range outs {
			if out.IsLockedWithKey(pubKeyHash) {
				accumulated += out.Value
				hashedUTxnoutIdxs[txnHash] = append(hashedUTxnoutIdxs[txnHash], outIdx)

				if accumulated >= amount {
					return false
				}
			}
		}
		return true
	})

	return accumulated, hashedUTxnoutIdxs
}
