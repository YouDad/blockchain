package core

import (
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/YouDad/blockchain/global"
	"github.com/YouDad/blockchain/global/mempool"
	"github.com/YouDad/blockchain/log"
	"github.com/YouDad/blockchain/types"
	"github.com/YouDad/blockchain/utils"
	"github.com/YouDad/blockchain/wallet"
)

// 交易地址: Unspent TxnOutputs
type UTXOSet struct {
	db    *global.UTXOSetDB
	bc    *Blockchain
	group int
}

func (set *UTXOSet) Clear() {
	set.db.Clear(set.group)
}

func (set *UTXOSet) Get(key interface{}) (value []byte) {
	return set.db.Get(set.group, key)
}

func (set *UTXOSet) Set(key interface{}, value []byte) {
	set.db.Set(set.group, key, value)
}

func (set *UTXOSet) Delete(key interface{}) {
	set.db.Delete(set.group, key)
}

func (set *UTXOSet) Foreach(fn func(k, v []byte) bool) {
	set.db.Foreach(set.group, fn)
}
func GetUTXOSet(group int) *UTXOSet {
	return &UTXOSet{global.GetUTXOSetDB(), GetBlockchain(group), group}
}

func (set *UTXOSet) Update(b *types.Block) {
	for _, txn := range b.Txns {
		if txn.IsCoinbase() == false {
			for _, vin := range txn.Vin {
				updatedOuts := []types.TxnOutput{}
				outsBytes := set.Get(vin.VoutHash)
				if len(outsBytes) == 0 {
					set.Reindex()
					set.bc.TxnReindex()
					return
				}
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
		set.Set(txn.Hash(), utils.Encode(newOutputs))
	}
}

func (set *UTXOSet) Reverse(b *types.Block) {
	for i := range b.Txns {
		txn := b.Txns[len(b.Txns)-i-1]
		if !txn.IsCoinbase() {
			set.Delete(txn.Hash())
			for _, vin := range txn.Vin {
				var txos []types.TxnOutput
				txosBytes := set.Get(vin.VoutHash)
				if len(txosBytes) != 0 {
					txos = BytesToTxnOutputs(txosBytes)
				}
				txos = append(txos, types.TxnOutput{
					Value:      vin.VoutValue,
					PubKeyHash: vin.PubKey.Hash(),
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

// 构造新的交易
func (set *UTXOSet) CreateTransaction(from, to string, amount int64) (*types.Transaction, error) {
	// 找到发送者的私钥
	wallets, err := wallet.GetWallets()
	if err != nil {
		return nil, err
	}

	fromWallet, have := wallets[from]
	if !have {
		return nil, errors.New(fmt.Sprintf("You haven't %s's PrivateKey", from))
	}

	// 用公钥找到一定数量的余额
	sum, utxos, values := set.findUTXOs(fromWallet.PublicKey, amount)

	if sum < amount {
		return nil, errors.New("Not enough BTC")
	}

	// 构造TxnInput
	ins := []types.TxnInput{}
	for txnHash, outIdxs := range utxos {
		txnHashByte, err := hex.DecodeString(txnHash)
		if err != nil {
			return nil, err
		}

		for i, outIdx := range outIdxs {
			ins = append(ins, types.TxnInput{
				VoutHash:  txnHashByte,
				VoutIndex: outIdx,
				VoutValue: values[txnHash][i],
				Signature: nil,
				PubKey:    fromWallet.PublicKey,
			})
		}
	}

	// 构造TxnOutput
	outs := []types.TxnOutput{*NewTxnOutput(to, amount)}
	if sum > amount {
		outs = append(outs, *NewTxnOutput(from, sum-amount))
	}

	// 交易签名
	txn := types.Transaction{Vin: ins, Vout: outs}
	err = set.bc.SignTransaction(&txn, fromWallet.PrivateKey)
	return &txn, err
}

func (set *UTXOSet) FindUTXOByHash(pubKey types.PublicKey) []types.TxnOutput {
	utxos := []types.TxnOutput{}

	set.Foreach(func(k, v []byte) bool {
		outs := BytesToTxnOutputs(v)

		for _, out := range outs {
			if out.IsLockedWithKey(pubKey) {
				utxos = append(utxos, out)
			}
		}
		return true
	})

	return utxos
}

// 用公钥找一定数量余额
func (set *UTXOSet) findUTXOs(pubKey types.PublicKey, amount int64) (int64, map[string][]int, map[string][]int64) {
	hashedUTXOIdxs := make(map[string][]int)
	hashedUTXOValues := make(map[string][]int64)
	var sum int64 = 0

	set.Foreach(func(k, v []byte) bool {
		txnOutputs := BytesToTxnOutputs(v)

		for txnOutputIndex, txnOutput := range txnOutputs {
			if txnOutput.IsLockedWithKey(pubKey) {
				outs, hashs, indexs := mempool.ExpandTxnOutput(set.group, txnOutput, k, txnOutputIndex)

				for i := range outs {
					str := hashs[i].String()

					sum += outs[i].Value
					hashedUTXOIdxs[str] = append(hashedUTXOIdxs[str], indexs[i])
					hashedUTXOValues[str] = append(hashedUTXOValues[str], outs[i].Value)

					if sum >= amount {
						return false
					}
				}
			}
		}
		return true
	})

	return sum, hashedUTXOIdxs, hashedUTXOValues
}
