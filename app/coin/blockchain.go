package coin

import (
	"encoding/hex"
	"log"
	"os"

	"github.com/YouDad/blockchain/app"
	"github.com/YouDad/blockchain/core"
)

type coinBlockchain struct {
	*core.Blockchain
}

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

	acc, validOutputs := bc.FindSpendableOutputs(from, amount)

	if acc < amount {
		log.Println("ERROR: Not enough funds")
		os.Exit(1)
	}

	// Build a list of inputs
	for txid, outs := range validOutputs {
		txID, err := hex.DecodeString(txid)
		if err != nil {
			log.Panic(err)
		}

		for _, out := range outs {
			input := TXInput{txID, out, from}
			inputs = append(inputs, input)
		}
	}

	// Build a list of outputs
	outputs = append(outputs, TXOutput{amount, to})
	if acc > amount {
		outputs = append(outputs, TXOutput{acc - amount, from}) // a change
	}

	tx := Transaction{nil, inputs, outputs}
	tx.SetID()

	return &tx
}

// FindUnspentTransactions returns a list of transactions containing unspent outputs
func (bc *coinBlockchain) FindUnspentTransactions(address string) []Transaction {
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

				if out.CanBeUnlockedWith(address) {
					unspentTXs = append(unspentTXs, *tx)
				}
			}

			if tx.IsCoinbase() == false {
				for _, in := range tx.Vin {
					if in.CanUnlockOutputWith(address) {
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
func (bc *coinBlockchain) FindUTXO(address string) []TXOutput {
	var UTXOs []TXOutput
	unspentTransactions := bc.FindUnspentTransactions(address)

	for _, tx := range unspentTransactions {
		for _, out := range tx.Vout {
			if out.CanBeUnlockedWith(address) {
				UTXOs = append(UTXOs, out)
			}
		}
	}

	return UTXOs
}

// FindSpendableOutputs finds and returns unspent outputs to reference in inputs
func (bc *coinBlockchain) FindSpendableOutputs(address string, amount int) (int, map[string][]int) {
	unspentOutputs := make(map[string][]int)
	unspentTXs := bc.FindUnspentTransactions(address)
	accumulated := 0

Work:
	for _, tx := range unspentTXs {
		txID := hex.EncodeToString(tx.ID)

		for outIdx, out := range tx.Vout {
			if out.CanBeUnlockedWith(address) {
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
