package core

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/hex"
	"errors"

	"github.com/YouDad/blockchain/app"
	"github.com/YouDad/blockchain/core"
	"github.com/YouDad/blockchain/log"
)

type CoinBlockchain struct {
	*core.Blockchain
}

const genesisBlockData = "The Times 03/Jan/2009 Chancellor on brink of second bailout for banks"

func NewBlockchain() *CoinBlockchain {
	return &CoinBlockchain{core.NewBlockchain()}
}

func CreateBlockchain(address string) {
	core.InitCore(core.Config{
		GetGenesis: func() app.App {
			return GetCoinApp([]*Transaction{
				NewCoinbaseTX(address, genesisBlockData),
			})
		},
	})
	core.CreateBlockchain().Close()
}

func NewProofOfWork(b *core.Block) *core.ProofOfWork {
	return core.NewProofOfWork(b)
}

func (bc *CoinBlockchain) MineBlock(Transactions []*Transaction) *core.Block {
	for _, tx := range Transactions {
		if !bc.VerifyTransaction(tx) {
			log.Panic("ERROR: Invalid transaction")
		}
	}
	return bc.Blockchain.MineBlock(GetCoinApp(Transactions))
}

// FindSpendableOutputs finds and returns unspent outputs to reference in inputs
func (bc *CoinBlockchain) FindSpendableOutputs(pubKeyHash []byte, amount int) (int, map[string][]int) {
	unspentOutputs := make(map[string][]int)
	unspentTXs := bc.FindUnspentTransactions(pubKeyHash)
	accumulated := 0

Work:
	for _, tx := range unspentTXs {
		txID := hex.EncodeToString(tx.ID)

		for outIdx, out := range tx.Vout {
			if out.IsLockedWithKey(pubKeyHash) {
				accumulated += out.Value
				unspentOutputs[txID] = append(unspentOutputs[txID], outIdx)
			}
			if accumulated >= amount {
				break Work
			}
		}
	}

	return accumulated, unspentOutputs
}

// FindUnspentTransactions returns a list of transactions containing unspent outputs
func (bc *CoinBlockchain) FindUnspentTransactions(pubKeyHash []byte) []Transaction {
	var unspentTXs []Transaction
	spentTXOs := make(map[string][]int)
	iter := bc.Begin()

	for {
		block := iter.Next()
		if block == nil {
			break
		}

		txs := block.App.(*CoinApp)

		for _, tx := range txs.Transactions {
			txID := hex.EncodeToString(tx.ID)

		Outputs:
			for outIdx, out := range tx.Vout {
				// Was the output spent?
				if spentTXOs[txID] != nil {
					for _, spentOut := range spentTXOs[txID] {
						if spentOut == outIdx {
							continue Outputs
						}
					}
				}

				if out.IsLockedWithKey(pubKeyHash) {
					unspentTXs = append(unspentTXs, *tx)
				}
			}

			if tx.IsCoinbase() == false {
				for _, in := range tx.Vin {
					if in.UsesKey(pubKeyHash) {
						inTxID := hex.EncodeToString(in.Txid)
						spentTXOs[inTxID] = append(spentTXOs[inTxID], in.Vout)
					}
				}
			}
		}
	}

	return unspentTXs
}

// FindUTXO finds and returns all unspent transaction outputs
func (bc *CoinBlockchain) FindUTXO() map[string]TXOutputs {
	UTXO := make(map[string]TXOutputs)
	spentTXOs := make(map[string][]int)
	iter := bc.Begin()

	for {
		block := iter.Next()
		if block == nil {
			break
		}

		txs := block.App.(*CoinApp)

		for _, tx := range txs.Transactions {
			txID := hex.EncodeToString(tx.ID)

		Outputs:
			for outIdx, out := range tx.Vout {
				// Was the output spent?
				if spentTXOs[txID] != nil {
					for _, spentOutIdx := range spentTXOs[txID] {
						if spentOutIdx == outIdx {
							continue Outputs
						}
					}
				}

				outs := UTXO[txID]
				outs.Outputs = append(outs.Outputs, out)
				UTXO[txID] = outs
			}

			if tx.IsCoinbase() == false {
				for _, in := range tx.Vin {
					inTxID := hex.EncodeToString(in.Txid)
					spentTXOs[inTxID] = append(spentTXOs[inTxID], in.Vout)
				}
			}
		}

		if len(block.PrevBlockHash) == 0 {
			break
		}
	}

	return UTXO
}

// FindAllUTXO finds all unspent transaction outputs and returns transactions with spent outputs removed
func (bc *CoinBlockchain) FindAllUTXO() map[string]TXOutputs {
	UTXO := make(map[string]TXOutputs)
	spentTXOs := make(map[string][]int)
	iter := bc.Begin()

	for {
		block := iter.Next()
		if block == nil {
			break
		}

		txs := block.App.(*CoinApp)

		for _, tx := range txs.Transactions {
			txID := hex.EncodeToString(tx.ID)

		Outputs:
			for outIdx, out := range tx.Vout {
				// Was the output spent?
				if spentTXOs[txID] != nil {
					for _, spentOutIdx := range spentTXOs[txID] {
						if spentOutIdx == outIdx {
							continue Outputs
						}
					}
				}

				outs := UTXO[txID]
				outs.Outputs = append(outs.Outputs, out)
				UTXO[txID] = outs
			}

			if tx.IsCoinbase() == false {
				for _, in := range tx.Vin {
					inTxID := hex.EncodeToString(in.Txid)
					spentTXOs[inTxID] = append(spentTXOs[inTxID], in.Vout)
				}
			}
		}
	}

	return UTXO
}

// FindTransaction finds a transaction by its ID
func (bc *CoinBlockchain) FindTransaction(ID []byte) (Transaction, error) {
	iter := bc.Begin()

	for {
		block := iter.Next()
		if block == nil {
			break
		}

		txs := block.App.(*CoinApp)

		for _, tx := range txs.Transactions {
			if bytes.Compare(tx.ID, ID) == 0 {
				return *tx, nil
			}
		}
	}

	return Transaction{}, errors.New("Transaction is not found")
}

// SignTransaction signs inputs of a Transaction
func (bc *CoinBlockchain) SignTransaction(tx *Transaction, privKey ecdsa.PrivateKey) {
	prevTXs := make(map[string]Transaction)

	for _, vin := range tx.Vin {
		prevTX, err := bc.FindTransaction(vin.Txid)
		if err != nil {
			log.Panic(err)
		}
		prevTXs[hex.EncodeToString(prevTX.ID)] = prevTX
	}

	tx.Sign(privKey, prevTXs)
}

// VerifyTransaction verifies transaction input signatures
func (bc *CoinBlockchain) VerifyTransaction(tx *Transaction) bool {
	if tx.IsCoinbase() {
		return true
	}

	prevTXs := make(map[string]Transaction)

	for _, vin := range tx.Vin {
		prevTX, err := bc.FindTransaction(vin.Txid)
		if err != nil {
			log.Panic(err)
		}
		prevTXs[hex.EncodeToString(prevTX.ID)] = prevTX
	}

	return tx.Verify(prevTXs)
}
