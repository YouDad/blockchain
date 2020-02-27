package core

import (
	"bytes"
	"encoding/hex"

	"github.com/YouDad/blockchain/conf"
	"github.com/YouDad/blockchain/log"
	"github.com/YouDad/blockchain/store"
	"github.com/YouDad/blockchain/types"
	"github.com/YouDad/blockchain/utils"
)

type Blockchain struct {
	store.IDatabase
}

type BlockchainIterator struct {
	*Blockchain
	next types.HashValue
}

func (bc *Blockchain) Begin() *BlockchainIterator {
	lastestBlock := BytesToBlock(bc.SetTable(conf.BLOCKS).Get("lastest"))
	return &BlockchainIterator{bc, lastestBlock.Hash()}
}

func (iter *BlockchainIterator) Next() (nextBlock *Block) {
	iter.SetTable(conf.BLOCKS)
	if iter.next == nil {
		return nil
	}
	nextBlock = BytesToBlock(iter.Get(iter.next))
	iter.next = nextBlock.PrevHash
	return nextBlock
}

func CreateBlockchain(minerAddress string) Blockchain {
	bc := Blockchain{store.CreateDatabase()}
	bc.SetTable(conf.BLOCKS).Clear()
	genesis := NewBlock(nil, 1, []*Transaction{NewCoinbaseTxn(minerAddress)})
	bytes := utils.Encode(genesis)
	bc.Set(genesis.Hash(), bytes)
	bc.Set(genesis.Height, bytes)
	bc.Set("genesis", bytes)
	bc.Set("lastest", bytes)
	return bc
}

func CreateBlockchainFromGenesis(b *Block) *Blockchain {
	log.NotImplement()
	return nil
}

func IsExists() bool {
	log.NotImplement()
	return true
}

func GetBlockchain() *Blockchain {
	return &Blockchain{store.GetDatabase()}
}

func (bc *Blockchain) GetGenesis() *Block {
	bc.SetTable(conf.BLOCKS)
	return BytesToBlock(bc.Get("genesis"))
}

func (bc *Blockchain) GetHeight() int {
	bc.SetTable(conf.BLOCKS)
	return BytesToBlock(bc.Get("genesis")).Height
}

func (bc *Blockchain) AddBlock(b *Block) {
	log.NotImplement()
}

func (bc *Blockchain) MineBlock(txns []*Transaction) *Block {
	for _, txn := range txns {
		if !bc.VerifyTransaction(txn) {
			log.Errln("Invalid transaction")
		}
	}

	lastestBlock := BytesToBlock(bc.SetTable(conf.BLOCKS).Get("lastest"))
	newBlock := NewBlock(lastestBlock.Hash(), lastestBlock.Height+1, txns)
	newBlockBytes := utils.Encode(newBlock)
	bc.SetTable(conf.BLOCKS)
	bc.Set("lastest", newBlockBytes)
	bc.Set(newBlock.Hash(), newBlockBytes)
	bc.Set(newBlock.Height, newBlockBytes)
	return newBlock
}

func (bc *Blockchain) FindUTXO() map[string][]TxnOutput {
	utxos := map[string][]TxnOutput{}
	// spent transaction outputs
	stxos := map[string]map[int]bool{}
	iter := bc.Begin()
	for {
		block := iter.Next()
		if block == nil {
			break
		}

		for _, txn := range block.Txns {
			hash := hex.EncodeToString(txn.Hash[:])

			for i, out := range txn.Vout {
				if stxos[hash] != nil {
					_, ok := stxos[hash][i]
					if ok {
						continue
					}
				}

				utxos[hash] = append(utxos[hash], out)
			}

			if !txn.IsCoinbase() {
				for _, in := range txn.Vin {
					inHash := hex.EncodeToString(in.VoutHash[:])
					stxos[inHash][in.VoutIndex] = true
				}
			}
		}

		if len(block.PrevHash) == 0 {
			break
		}
	}
	return utxos
}

func (bc *Blockchain) FindTxn(hash types.HashValue) Transaction {
	iter := bc.Begin()

	for {
		block := iter.Next()
		if block == nil {
			break
		}

		for _, txn := range block.Txns {
			if bytes.Compare(txn.Hash, hash) == 0 {
				return *txn
			}
		}
	}

	log.Errln("Transaction is not found")
	return Transaction{}
}

func (bc *Blockchain) SignTransaction(txn *Transaction, sk types.PrivateKey) {
	hashedTxn := make(map[string]Transaction)

	for _, vin := range txn.Vin {
		vinTxn := bc.FindTxn(vin.VoutHash)
		hashedTxn[hex.EncodeToString(vinTxn.Hash)] = vinTxn
	}

	txn.Sign(sk, hashedTxn)
}

func (bc *Blockchain) VerifyTransaction(txn *Transaction) bool {
	if txn.IsCoinbase() {
		return true
	}

	prevTXs := make(map[string]Transaction)

	for _, vin := range txn.Vin {
		prevTX := bc.FindTxn(vin.VoutHash)
		prevTXs[hex.EncodeToString(prevTX.Hash)] = prevTX
	}

	return txn.Verify(prevTXs)
}
