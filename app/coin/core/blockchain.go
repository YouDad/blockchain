package core

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"log"

	"github.com/YouDad/blockchain/app"
	"github.com/YouDad/blockchain/app/coin/wallet"
	"github.com/YouDad/blockchain/core"
)

type coinBlockchain struct {
	*core.Blockchain
}

const genesisBlockData = "The Times 03/Jan/2009 Chancellor on brink of second bailout for banks"

func NewBlockchain() *coinBlockchain {
	return &coinBlockchain{core.NewBlockchain()}
}

func CreateBlockchain(address string) {
	core.InitCore(core.CoreConfig{
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

// NewUTXOTransaction creates a new transaction
func (bc *coinBlockchain) NewUTXOTransaction(from, to string, amount int) *Transaction {
	var inputs []TXInput
	var outputs []TXOutput

	wallets, err := wallet.NewWallets()
	if err != nil {
		log.Panic(err)
	}

	srcWallet := wallets.GetWallet(from)
	pubKeyHash := wallet.HashPubKey(srcWallet.PublicKey)
	acc, validOutputs := bc.FindSpendableOutputs(pubKeyHash, amount)

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
	bc.SignTransaction(&tx, srcWallet.PrivateKey)

	return &tx
}

// FindUnspentTransactions returns a list of transactions containing unspent outputs
func (bc *coinBlockchain) FindUnspentTransactions(pubKeyHash []byte) []Transaction {
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
func (bc *coinBlockchain) FindUTXO(pubKeyHash []byte) []TXOutput {
	var UTXOs []TXOutput
	unspentTransactions := bc.FindUnspentTransactions(pubKeyHash)

	for _, tx := range unspentTransactions {
		for _, out := range tx.Vout {
			if out.IsLockedWithKey(pubKeyHash) {
				UTXOs = append(UTXOs, out)
			}
		}
	}

	return UTXOs
}

// FindSpendableOutputs finds and returns unspent outputs to reference in inputs
func (bc *coinBlockchain) FindSpendableOutputs(pubKeyHash []byte, amount int) (int, map[string][]int) {
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

// FindTransaction finds a transaction by its ID
func (bc *coinBlockchain) FindTransaction(ID []byte) (Transaction, error) {
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
func (bc *coinBlockchain) SignTransaction(tx *Transaction, privKey ecdsa.PrivateKey) {
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
func (bc *coinBlockchain) VerifyTransaction(tx *Transaction) bool {
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
