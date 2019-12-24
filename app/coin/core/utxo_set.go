package core

import (
	"encoding/hex"

	"github.com/YouDad/blockchain/app/coin/wallet"
	"github.com/YouDad/blockchain/core"
	"github.com/YouDad/blockchain/log"
)

// UTXOSet represents UTXO set
type UTXOSet struct {
	*CoinBlockchain
}

func NewUTXOSet() *UTXOSet {
	return &UTXOSet{NewBlockchain()}
}

// NewUTXOTransaction creates a new transaction
func (u *UTXOSet) NewUTXOTransaction(from, to string, amount int) *Transaction {
	var inputs []TXInput
	var outputs []TXOutput

	wallets, err := wallet.NewWallets(core.CoreConfig.WalletFile)
	if err != nil {
		log.Panic(err)
	}

	srcWallet := wallets.GetWallet(from)
	pubKeyHash := wallet.HashPubKey(srcWallet.PublicKey)
	acc, validOutputs := u.FindSpendableOutputs(pubKeyHash, amount)

	if acc < amount {
		log.Panic("ERROR: Not enough funds")
	}

	// Build a list of inputs
	for txid, outs := range validOutputs {
		txID, err := hex.DecodeString(txid)
		if err != nil {
			log.Panic(err)
		}

		for _, out := range outs {
			input := TXInput{txID, out, nil, srcWallet.PublicKey}
			inputs = append(inputs, input)
		}
	}

	// Build a list of outputs
	outputs = append(outputs, *NewTXOutput(to, amount))
	if acc > amount {
		outputs = append(outputs, *NewTXOutput(from, acc-amount)) // a change
	}

	tx := Transaction{nil, inputs, outputs}
	tx.ID = tx.Hash()
	u.SignTransaction(&tx, srcWallet.PrivateKey)

	return &tx
}

// FindSpendableOutputs finds and returns unspent outputs to reference in inputs
func (u UTXOSet) FindSpendableOutputs(pubkeyHash []byte, amount int) (int, map[string][]int) {
	unspentOutputs := make(map[string][]int)
	accumulated := 0

	err := u.UTXOSet().Foreach(func(k, v []byte) bool {
		txID := hex.EncodeToString(k)
		outs := DeserializeOutputs(v)

		for outIdx, out := range outs.Outputs {
			if out.IsLockedWithKey(pubkeyHash) {
				accumulated += out.Value
				unspentOutputs[txID] = append(unspentOutputs[txID], outIdx)

				if accumulated >= amount {
					return false
				}
			}
		}
		return true
	})
	if err != nil {
		log.Panic(err)
	}

	return accumulated, unspentOutputs
}

// FindUTXO finds UTXO for a public key hash
func (u UTXOSet) FindUTXO(pubKeyHash []byte) []TXOutput {
	var UTXOs []TXOutput

	err := u.UTXOSet().Foreach(func(k, v []byte) (isContinue bool) {
		outs := DeserializeOutputs(v)

		for _, out := range outs.Outputs {
			if out.IsLockedWithKey(pubKeyHash) {
				UTXOs = append(UTXOs, out)
			}
		}
		return true
	})
	if err != nil {
		log.Panic(err)
	}

	return UTXOs
}

// CountTransactions returns the number of transactions in the UTXO set
func (u UTXOSet) CountTransactions() int {
	counter := 0

	err := u.UTXOSet().Foreach(func(k, v []byte) (isContinue bool) {
		counter++
		return true
	})
	if err != nil {
		log.Panic(err)
	}

	return counter
}

// Reindex rebuilds the UTXO set
func (u UTXOSet) Reindex() {
	u.UTXOSet().Clear()
	UTXO := u.CoinBlockchain.FindUTXO()

	for txID, outs := range UTXO {
		key, err := hex.DecodeString(txID)
		if err != nil {
			log.Panic(err)
		}

		u.UTXOSet().Set(key, outs.Serialize())
	}
}

// Update updates the UTXO set with transactions from the Block
// The Block is considered to be the tip of a blockchain
func (u UTXOSet) Update(block *core.Block) {
	txs := block.App.(*CoinApp)

	for _, tx := range txs.Transactions {
		if tx.IsCoinbase() == false {
			for _, vin := range tx.Vin {
				updatedOuts := TXOutputs{}
				outsBytes := u.UTXOSet().Get(vin.Txid)
				outs := DeserializeOutputs(outsBytes)

				for outIdx, out := range outs.Outputs {
					if outIdx != vin.Vout {
						updatedOuts.Outputs = append(updatedOuts.Outputs, out)
					}
				}

				if len(updatedOuts.Outputs) == 0 {
					err := u.UTXOSet().Delete(vin.Txid)
					if err != nil {
						log.Panic(err)
					}
				} else {
					u.UTXOSet().Set(vin.Txid, updatedOuts.Serialize())
				}
			}
		}

		newOutputs := TXOutputs{}
		for _, out := range tx.Vout {
			newOutputs.Outputs = append(newOutputs.Outputs, out)
		}

		u.UTXOSet().Set(tx.ID, newOutputs.Serialize())
	}
}
