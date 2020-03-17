package core

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/YouDad/blockchain/global"
	"github.com/YouDad/blockchain/log"
	"github.com/YouDad/blockchain/types"
	"github.com/YouDad/blockchain/utils"
)

type Blockchain struct {
	*global.BlocksDB
}

type BlockchainIterator struct {
	*Blockchain
	next types.HashValue
}

func (bc *Blockchain) Begin() *BlockchainIterator {
	return &BlockchainIterator{bc, bc.GetLastest().Hash}
}

func (iter *BlockchainIterator) Next() (nextBlock *types.Block) {
	if iter.next == nil {
		return nil
	}
	nextBlock = BytesToBlock(global.GetBlocksDB().Get(iter.next))
	if nextBlock != nil {
		iter.next = nextBlock.PrevHash
	}
	return nextBlock
}

func CreateBlockchain(minerAddress string) *Blockchain {
	global.CreateDatabase()
	bc := GetBlockchain()
	bc.Clear()
	global.SetHeight(-1)
	genesis := NewBlock(nil, 1, 0, []*types.Transaction{NewCoinbaseTxn(minerAddress)})
	bytes := utils.Encode(genesis)
	bc.Set(genesis.Hash, bytes)
	bc.Set(genesis.Height, bytes)
	bc.Set("genesis", bytes)
	bc.Set("lastest", bytes)
	GetUTXOSet().Reindex()
	return bc
}

func CreateBlockchainFromGenesis(genesis *types.Block) *Blockchain {
	global.CreateDatabase()
	bc := GetBlockchain()
	bc.Clear()
	bytes := utils.Encode(genesis)
	bc.Set(genesis.Hash, bytes)
	bc.Set(genesis.Height, bytes)
	bc.Set("genesis", bytes)
	bc.Set("lastest", bytes)
	GetUTXOSet().Reindex()
	return bc
}

func GetBlockchain() *Blockchain {
	return &Blockchain{global.GetBlocksDB()}
}

func (bc *Blockchain) GetGenesis() *types.Block {
	return BytesToBlock(bc.Get("genesis"))
}

func (bc *Blockchain) GetLastest() *types.Block {
	lastest := BytesToBlock(bc.Get("lastest"))
	global.SetHeight(lastest.Height)
	return lastest
}

func (bc *Blockchain) GetHeight() int32 {
	return bc.GetLastest().Height
}

func (bc *Blockchain) AddBlock(b *types.Block) {
	if b == nil {
		return
	}

	if bc.Get(b.Hash) != nil {
		return
	}

	bytes := utils.Encode(b)
	if global.GetHeight() < b.Height {
		bc.Set("lastest", bytes)
		bc.Set(b.Hash, bytes)
		bc.Set(b.Height, bytes)
		global.SetHeight(b.Height)
	}
}

func (bc *Blockchain) MineBlock(txns []*types.Transaction) *types.Block {
	for _, txn := range txns {
		if !bc.VerifyTransaction(*txn) {
			log.Errln("Invalid transaction")
		}
	}

	lastest := bc.GetLastest()
	difficulty := lastest.Difficulty
	height := lastest.Height + 1
	if height%60 == 0 {
		lastDiff := BytesToBlock(bc.Get(height - 60))
		thisDiff := BytesToBlock(bc.Get(height - 1))
		difficulty *= 59 * 60 * 1e9 / float64(thisDiff.Timestamp-lastDiff.Timestamp)
	}

	newBlock := NewBlock(lastest.Hash, difficulty, height, txns)
	if newBlock == nil {
		return nil
	}
	log.Infoln("New Block", lastest.Hash, difficulty, height)
	bc.AddBlock(newBlock)
	global.SyncMutex.Unlock()
	return newBlock
}

func (bc *Blockchain) FindUTXO() map[string][]types.TxnOutput {
	utxos := make(map[string][]types.TxnOutput)
	// spent transaction outputs
	stxos := make(map[string]map[int]bool)
	iter := bc.Begin()
	for {
		block := iter.Next()
		if block == nil {
			break
		}

		for _, txn := range block.Txns {
			hash := hex.EncodeToString(txn.Hash[:]) // 交易哈希

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
					_, ok := stxos[inHash]
					if !ok {
						stxos[inHash] = make(map[int]bool)
					}
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

func (bc *Blockchain) FindTxn(hash types.HashValue) (types.Transaction, error) {
	var block *types.Block
	iter := bc.Begin()

	for block = iter.Next(); block != nil; block = iter.Next() {
		for _, txn := range block.Txns {
			if bytes.Compare(txn.Hash, hash) == 0 {
				return *txn, nil
			}
		}
	}

	return types.Transaction{}, errors.New(fmt.Sprintf("Transaction is not found, %x", hash))
}

func (bc *Blockchain) SignTransaction(txn *types.Transaction, sk types.PrivateKey) error {
	hashedTxn := make(map[string]types.Transaction)

	for _, vin := range txn.Vin {
		var vinTxn types.Transaction
		vinTxn, err := bc.FindTxn(vin.VoutHash)
		if err != nil {
			return err
		}
		hashedTxn[hex.EncodeToString(vinTxn.Hash)] = vinTxn
	}

	txn.Sign(sk, hashedTxn)
	return nil
}

func (bc *Blockchain) VerifyTransaction(txn types.Transaction) bool {
	if txn.IsCoinbase() {
		return true
	}

	prevTXs := make(map[string]types.Transaction)

	for _, vin := range txn.Vin {
		prevTX, err := bc.FindTxn(vin.VoutHash)
		if err != nil {
			return false
		}
		prevTXs[hex.EncodeToString(prevTX.Hash)] = prevTX
	}

	return txn.Verify(prevTXs)
}
