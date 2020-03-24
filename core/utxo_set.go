package core

import (
	"encoding/hex"

	"github.com/YouDad/blockchain/global"
	"github.com/YouDad/blockchain/log"
	"github.com/YouDad/blockchain/types"
	"github.com/YouDad/blockchain/utils"
	"github.com/YouDad/blockchain/wallet"
)

type UTXOSet struct {
	*global.UTXOSetDB
	bc *Blockchain
}

func GetUTXOSet() *UTXOSet {
	return &UTXOSet{global.GetUTXOSetDB(), GetBlockchain()}
}

func (set *UTXOSet) Update(b *types.Block) {
	for _, txn := range b.Txns {
		if txn.IsCoinbase() == false {
			for _, vin := range txn.Vin {
				updatedOuts := []types.TxnOutput{}
				outsBytes := set.Get(vin.VoutHash)
				outs := BytesToTxnOutputs(outsBytes)

				for outIdx, out := range outs {
					if outIdx != vin.VoutIndex {
						updatedOuts = append(updatedOuts, out)
					}
				}

				if len(updatedOuts) == 0 {
					set.Delete(vin.VoutHash)
				} else {
					set.Set(vin.VoutHash, utils.Encode(updatedOuts))
				}
			}
		}

		newOutputs := []types.TxnOutput{}
		for _, out := range txn.Vout {
			newOutputs = append(newOutputs, out)
		}
		set.Set(txn.Hash, utils.Encode(newOutputs))
	}
}

func (set *UTXOSet) Reverse(b *types.Block) {
	for _, txn := range b.Txns {
		if !txn.IsCoinbase() {
			set.Delete(txn.Hash)
			for _, vin := range txn.Vin {
				var txos []types.TxnOutput
				txosBytes := set.Get(vin.VoutHash)
				if len(txosBytes) != 0 {
					txos = BytesToTxnOutputs(txosBytes)
				}
				txos = append(txos, types.TxnOutput{
					Value:      vin.VoutValue,
					PubKeyHash: vin.PubKeyHash,
				})
				set.Set(vin.VoutHash, utils.Encode(txos))
			}
		}
	}
}

func (set *UTXOSet) Reindex() {
	hashedUtxos := set.bc.FindUTXO()
	set.Clear()

	for txnHash, utxos := range hashedUtxos {
		hash, err := hex.DecodeString(txnHash)
		log.Err(err)
		set.Set(hash, utils.Encode(utxos))
	}
}

func (set *UTXOSet) NewUTXOTransaction(from, to string, amount int64) (*types.Transaction, error) {
	var ins []types.TxnInput
	var outs []types.TxnOutput

	wallets, err := wallet.NewWallets()
	log.Err(err)

	srcWallet, have := wallets.GetWallet(from)
	if !have {
		log.Errf("You haven't %s's PrivateKey", from)
	}
	pubKeyHash := wallet.HashPubKey(srcWallet.PublicKey)
	acc, utxos, values := set.FindUTXOs(pubKeyHash, amount)

	if acc < amount {
		log.Errln("Not enough BTC")
	}

	for txnHash, outIdxs := range utxos {
		txnHashByte, err := hex.DecodeString(txnHash)
		log.Err(err)

		for i, outIdx := range outIdxs {
			ins = append(ins, types.TxnInput{
				VoutHash:   txnHashByte,
				VoutIndex:  outIdx,
				VoutValue:  values[txnHash][i],
				Signature:  nil,
				PubKeyHash: srcWallet.PublicKey,
			})
		}
	}

	outs = append(outs, *NewTxnOutput(to, amount))
	if acc > amount {
		outs = append(outs, *NewTxnOutput(from, acc-amount))
	}

	txn := types.Transaction{
		Vin:  ins,
		Vout: outs,
	}
	err = set.bc.SignTransaction(&txn, srcWallet.PrivateKey)
	return &txn, err
}

func (set *UTXOSet) FindUTXOByHash(pubKeyHash []byte) []types.TxnOutput {
	utxos := []types.TxnOutput{}

	set.Foreach(func(k, v []byte) bool {
		outs := BytesToTxnOutputs(v)

		for _, out := range outs {
			// log.Tracef("%x %x\n", pubKeyHash, out.PubKeyHash)
			if out.IsLockedWithKey(pubKeyHash) {
				utxos = append(utxos, out)
			}
		}
		return true
	})

	return utxos
}

// 用公钥找一定数量余额
func (set *UTXOSet) FindUTXOs(pubKeyHash types.HashValue, amount int64) (int64, map[string][]int, map[string][]int64) {
	hashedUTXOIdxs := make(map[string][]int)
	hashedUTXOValues := make(map[string][]int64)
	var accumulated int64 = 0

	set.Foreach(func(k, v []byte) bool {
		txnHash := hex.EncodeToString(k)
		outs := BytesToTxnOutputs(v)

		for outIdx, out := range outs {
			if out.IsLockedWithKey(pubKeyHash) {
				accumulated += out.Value
				hashedUTXOIdxs[txnHash] = append(hashedUTXOIdxs[txnHash], outIdx)
				hashedUTXOValues[txnHash] = append(hashedUTXOValues[txnHash], out.Value)

				if accumulated >= amount {
					return false
				}
			}
		}
		return true
	})

	return accumulated, hashedUTXOIdxs, hashedUTXOValues
}
