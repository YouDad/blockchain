package core

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"log"
)

type CoinBlockchain struct {
	*Blockchain
}

const genesisBlockData = "The Times 03/Jan/2009 Chancellor on brink of second bailout for banks"

func NewBlockchain() *CoinBlockchain {
	return &CoinBlockchain{NewBlockchai()}
}

func CreateBlockchain(address string) {
	InitCore(Config{
		GetGenesis: func() CoinApp {
			return *GetCoinApp([]*Transaction{
				NewCoinbaseTX(address, genesisBlockData),
			})
		},
	})
	CreateBlockchai().Close()
}

func (bc *CoinBlockchain) MineBlock(Transactions []*Transaction) *Block {
	for _, tx := range Transactions {
		if !bc.VerifyTransaction(tx) {
			log.Panic("ERROR: Invalid transaction")
		}
	}
	return bc.Blockchain.MineBlock(*GetCoinApp(Transactions))
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

		txs := block.App

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

// FindTransaction finds a transaction by its ID
func (bc *CoinBlockchain) FindTransaction(ID []byte) (Transaction, error) {
	iter := bc.Begin()

	for {
		block := iter.Next()
		if block == nil {
			break
		}

		txs := block.App

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
