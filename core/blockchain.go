package core

import (
	"bytes"
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
	bc    *Blockchain
	group int
	next  types.HashValue
}

func (bc *Blockchain) Begin(group int) *BlockchainIterator {
	return &BlockchainIterator{bc, group, bc.GetLastest(group).Hash()}
}

func (iter *BlockchainIterator) Next() (nextBlock *types.Block) {
	if iter.next == nil {
		return nil
	}
	nextBlock = BytesToBlock(iter.bc.Get(iter.group, iter.next))
	if nextBlock != nil {
		iter.next = nextBlock.PrevHash
	}
	return nextBlock
}

func CreateBlockchain(minerAddress string) *Blockchain {
	group := global.GetGroup()
	bc := GetBlockchain()
	bc.Clear(group)
	genesis := NewBlock(group, nil, 1, 0, []*types.Transaction{NewCoinbaseTxn(minerAddress)})
	bytes := utils.Encode(genesis)
	bc.Set(group, genesis.Hash(), bytes)
	bc.Set(group, genesis.Height, bytes)
	bc.Set(group, "genesis", bytes)
	bc.Set(group, "lastest", bytes)
	GetUTXOSet().Reindex(group)
	return bc
}

func GetBlockchain() *Blockchain {
	return &Blockchain{global.GetBlocksDB()}
}

func (bc *Blockchain) GetGenesis(group int) *types.Block {
	return BytesToBlock(bc.Get(group, "genesis"))
}

func (bc *Blockchain) GetLastest(group int) *types.Block {
	return BytesToBlock(bc.Get(group, "lastest"))
}

func (bc *Blockchain) GetHeight(group int) int32 {
	lastest := bc.GetLastest(group)
	if lastest == nil {
		return -1
	}
	return lastest.Height
}

func (bc *Blockchain) AddBlock(group int, b *types.Block) {
	if b == nil {
		return
	}

	if bc.Get(group, b.Hash()) != nil {
		return
	}

	bytes := utils.Encode(b)
	if bc.GetHeight(group) < b.Height {
		if b.Height == 0 {
			bc.Set(group, "genesis", bytes)
		}
		bc.Set(group, "lastest", bytes)
		bc.Set(group, b.Hash(), bytes)
		bc.Set(group, b.Height, bytes)
	}
}

func (bc *Blockchain) MineBlock(group int, txns []*types.Transaction) *types.Block {
	for _, txn := range txns {
		if !bc.VerifyTransaction(group, *txn) {
			log.Errln("Invalid transaction")
		}
	}

	lastest := bc.GetLastest(group)
	difficulty := lastest.Difficulty
	height := lastest.Height + 1
	if height%60 == 0 {
		lastDiff := BytesToBlock(bc.Get(group, height-60))
		thisDiff := BytesToBlock(bc.Get(group, height-1))
		difficulty *= 59 * 60 * 1e9 / float64(thisDiff.Timestamp-lastDiff.Timestamp)
	}

	newBlock := NewBlock(group, lastest.Hash(), difficulty, height, txns)
	if newBlock == nil {
		return nil
	}
	log.Infof("NewBlock[%d]{%.2f} %s", height, difficulty, lastest.Hash())
	bc.AddBlock(group, newBlock)
	global.SyncMutex.Unlock()
	return newBlock
}

func (bc *Blockchain) FindUTXO(group int) map[string][]types.TxnOutput {
	utxos := make(map[string][]types.TxnOutput)
	// spent transaction outputs
	stxos := make(map[string]map[int]bool)
	iter := bc.Begin(group)
	for {
		block := iter.Next()
		if block == nil {
			break
		}

		for _, txn := range block.Txns {
			hash := txn.Hash().String()

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
					inHash := in.VoutHash.String()
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

func (bc *Blockchain) FindTxn(group int, hash types.HashValue) (types.Transaction, error) {
	var block *types.Block
	iter := bc.Begin(group)

	for block = iter.Next(); block != nil; block = iter.Next() {
		for _, txn := range block.Txns {
			if bytes.Compare(txn.Hash(), hash) == 0 {
				return *txn, nil
			}
		}
	}

	return types.Transaction{}, errors.New(fmt.Sprintf("Transaction is not found, %s", hash))
}

func (bc *Blockchain) SignTransaction(group int, txn *types.Transaction, sk types.PrivateKey) error {
	hashedTxn := make(map[string]types.Transaction)

	for _, vin := range txn.Vin {
		var vinTxn types.Transaction
		vinTxn, err := bc.FindTxn(group, vin.VoutHash)
		if err != nil {
			return err
		}
		hashedTxn[vinTxn.Hash().String()] = vinTxn
	}

	txn.Sign(sk, hashedTxn)
	return nil
}

func (bc *Blockchain) VerifyTransaction(group int, txn types.Transaction) bool {
	if txn.IsCoinbase() {
		return true
	}

	prevTXs := make(map[string]types.Transaction)

	for _, vin := range txn.Vin {
		prevTX, err := bc.FindTxn(group, vin.VoutHash)
		if err != nil {
			return false
		}
		prevTXs[prevTX.Hash().String()] = prevTX
	}

	return txn.Verify(prevTXs)
}
