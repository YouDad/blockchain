package core

import (
	"bytes"
	"encoding/hex"
	"time"

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
	return &BlockchainIterator{bc, lastestBlock.Hash}
}

func (iter *BlockchainIterator) Next() (nextBlock *Block) {
	if iter.next == nil {
		return nil
	}
	nextBlock = BytesToBlock(iter.SetTable(conf.BLOCKS).Get(iter.next))
	if nextBlock != nil {
		iter.next = nextBlock.PrevHash
	}
	return nextBlock
}

func CreateBlockchain(minerAddress string) *Blockchain {
	bc := Blockchain{store.CreateDatabase()}
	bc.SetTable(conf.BLOCKS).Clear()
	genesis := NewBlock(nil, 1, 0, []*Transaction{NewCoinbaseTxn(minerAddress)})
	bytes := utils.Encode(genesis)
	bc.SetTable(conf.BLOCKS).Set(genesis.Hash, bytes)
	bc.SetTable(conf.BLOCKS).Set(genesis.Height, bytes)
	bc.SetTable(conf.BLOCKS).Set("genesis", bytes)
	bc.SetTable(conf.BLOCKS).Set("lastest", bytes)
	return &bc
}

func CreateBlockchainFromGenesis(b *Block) *Blockchain {
	bc := Blockchain{store.CreateDatabase()}
	bc.SetTable(conf.BLOCKS).Clear()
	bytes := utils.Encode(b)
	bc.SetTable(conf.BLOCKS).Set(b.Hash, bytes)
	bc.SetTable(conf.BLOCKS).Set(b.Height, bytes)
	bc.SetTable(conf.BLOCKS).Set("genesis", bytes)
	bc.SetTable(conf.BLOCKS).Set("lastest", bytes)
	return &bc
}

func IsExists() bool {
	log.NotImplement()
	return true
}

func GetBlockchain() *Blockchain {
	return &Blockchain{store.GetDatabase()}
}

func (bc *Blockchain) GetGenesis() *Block {
	return BytesToBlock(bc.SetTable(conf.BLOCKS).Get("genesis"))
}

func (bc *Blockchain) GetLastest() *Block {
	return BytesToBlock(bc.SetTable(conf.BLOCKS).Get("lastest"))
}

func (bc *Blockchain) GetHeight() int32 {
	return BytesToBlock(bc.SetTable(conf.BLOCKS).Get("lastest")).Height
}

func (bc *Blockchain) AddBlock(b *Block) {
	if b == nil {
		return
	}
	if bc.SetTable(conf.BLOCKS).Get(b.Hash) != nil {
		return
	}

	lastestBlockBytes := bc.SetTable(conf.BLOCKS).Get("lastest")
	lastestBlock := BytesToBlock(lastestBlockBytes)
	bytes := utils.Encode(b)
	if lastestBlock.Height < b.Height {
		bc.SetTable(conf.BLOCKS).Set("lastest", bytes)
		bc.SetTable(conf.BLOCKS).Set(b.Hash, bytes)
		bc.SetTable(conf.BLOCKS).Set(b.Height, bytes)
	}
}

func (bc *Blockchain) MineBlock(txns []*Transaction) *Block {
	for _, txn := range txns {
		if !bc.VerifyTransaction(*txn) {
			log.Errln("Invalid transaction")
		}
	}

	lastest := BytesToBlock(bc.SetTable(conf.BLOCKS).Get("lastest"))
	difficulty := lastest.Difficulty
	height := lastest.Height + 1
	if height%60 == 0 {
		lastDiff := BytesToBlock(bc.SetTable(conf.BLOCKS).Get(height - 60))
		difficulty *= 3600.0 * 1e9 / float64(time.Now().UnixNano()-lastDiff.Timestamp)
	}

	newBlock := NewBlock(lastest.Hash, difficulty, height, txns)
	log.Infoln("NewBlock", lastest.Hash, difficulty, height)
	newBlockBytes := utils.Encode(newBlock)
	bc.SetTable(conf.BLOCKS).Set("lastest", newBlockBytes)
	bc.SetTable(conf.BLOCKS).Set(newBlock.Hash, newBlockBytes)
	bc.SetTable(conf.BLOCKS).Set(newBlock.Height, newBlockBytes)
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
	var block *Block
	iter := bc.Begin()

	for block = iter.Next(); block != nil; block = iter.Next() {
		for _, txn := range block.Txns {
			if bytes.Compare(txn.Hash, hash) == 0 {
				return *txn
			}
		}
	}

	log.Errln("Transaction is not found", hash)
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

func (bc *Blockchain) VerifyTransaction(txn Transaction) bool {
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
