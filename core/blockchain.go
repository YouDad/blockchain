package core

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"log"
	"os"

	"github.com/YouDad/blockchain/utils"
)

type Blockchain struct {
	*utils.Database
}

type BlockchainIterator struct {
	*Blockchain
	next []byte
}

func (bc *Blockchain) Begin() (iter *BlockchainIterator) {
	lastBlock := DeserializeBlock(bc.Blocks().GetLastest())
	return &BlockchainIterator{bc, lastBlock.Hash}
}

func (iter *BlockchainIterator) Next() (nextBlock *Block) {
	iter.Blockchain.Blocks()
	blockBytes := iter.Get(iter.next)
	if len(blockBytes) == 0 {
		return nil
	}
	nextBlock = DeserializeBlock(blockBytes)
	iter.next = nextBlock.PrevBlockHash
	return nextBlock
}

func (bc *Blockchain) mineBlock(data CoinApp) *Block {
	lastestBlock := DeserializeBlock(bc.GetLastest())
	newBlock := NewBlock(data, lastestBlock.Hash, lastestBlock.Height+1)
	bc.SetLastest(newBlock.Hash, newBlock.Serialize())
	bc.SetByInt(newBlock.Height, newBlock.Serialize())
	return newBlock
}

func (bc *Blockchain) AddBlock(block *Block) {
	if block == nil {
		return
	}
	if bc.Blocks().Get(block.Hash) != nil {
		return
	}

	lastestBlock := DeserializeBlock(bc.GetLastest())
	if lastestBlock.Height < block.Height {
		bc.SetLastest(block.Hash, block.Serialize())
		bc.SetByInt(block.Height, block.Serialize())
	}
}

func IsBlockchainExists() bool {
	return utils.IsDatabaseExists(CoreConfig.DatabaseFile)
}

func NewBlockchain() *Blockchain {
	if !utils.IsDatabaseExists(CoreConfig.DatabaseFile) {
		log.Println("No existing blockchain found. Create one to continue.")
		os.Exit(1)
	}

	return &Blockchain{utils.OpenDatabase(CoreConfig.DatabaseFile)}
}

func CreateBlockchainFromGenesis(genesis *Block) *Blockchain {
	if utils.IsDatabaseExists(CoreConfig.DatabaseFile) {
		log.Println("Blockchain existed, Create failed.")
		os.Exit(1)
	}

	db := utils.OpenDatabase(CoreConfig.DatabaseFile)
	db.Blocks().Clear()
	db.SetGenesis(genesis.Hash, genesis.Serialize())
	db.SetByInt(genesis.Height, genesis.Serialize())
	return &Blockchain{db}
}

func (bc *Blockchain) GetBestHeight() int {
	return DeserializeBlock(bc.GetLastest()).Height
}

func (bc *Blockchain) GetBlock(hash []byte) (*Block, error) {
	value := bc.Blocks().Get(hash)
	if len(value) == 0 {
		return nil, errors.New("Block is not found.")
	}
	return DeserializeBlock(value), nil
}

func (bc *Blockchain) GetBlockHashes() (hashes [][]byte) {
	iter := bc.Begin()
	for {
		block := iter.Next()
		if block == nil {
			break
		}
		hashes = append(hashes, block.Hash)
	}
	return hashes
}

const genesisBlockData = "The Times 03/Jan/2009 Chancellor on brink of second bailout for banks"

func CreateBlockchain(address string) {
	InitCore(Config{
		GetGenesis: func() CoinApp {
			return *GetCoinApp([]*Transaction{
				NewCoinbaseTX(address, genesisBlockData),
			})
		},
	})
	if utils.IsDatabaseExists(CoreConfig.DatabaseFile) {
		log.Panicln("Blockchain existed, Create failed.")
	}

	db := utils.OpenDatabase(CoreConfig.DatabaseFile)
	db.Blocks().Clear()
	genesis := NewBlock(CoreConfig.GetGenesis(), make([]byte, 32), 1)
	db.SetGenesis(genesis.Hash, genesis.Serialize())
	db.SetByInt(genesis.Height, genesis.Serialize())
	db.Close()
}

func (bc *Blockchain) MineBlock(Transactions []*Transaction) *Block {
	for _, tx := range Transactions {
		if !bc.VerifyTransaction(tx) {
			log.Panic("ERROR: Invalid transaction")
		}
	}
	return bc.mineBlock(*GetCoinApp(Transactions))
}

// FindUTXO finds and returns all unspent transaction outputs
func (bc *Blockchain) FindUTXO() map[string]TXOutputs {
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
func (bc *Blockchain) FindTransaction(ID []byte) (Transaction, error) {
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
func (bc *Blockchain) SignTransaction(tx *Transaction, privKey ecdsa.PrivateKey) {
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
func (bc *Blockchain) VerifyTransaction(tx *Transaction) bool {
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
