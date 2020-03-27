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
	db    *global.BlocksDB
	group int
}

func (bc *Blockchain) Clear() {
	bc.db.Clear(bc.group)
}

func (bc *Blockchain) Get(key interface{}) (value []byte) {
	return bc.db.Get(bc.group, key)
}

func (bc *Blockchain) Set(key interface{}, value []byte) {
	bc.db.Set(bc.group, key, value)
}

func (bc *Blockchain) Delete(key interface{}) {
	bc.db.Delete(bc.group, key)
}

func (bc *Blockchain) Foreach(fn func(k, v []byte) bool) {
	bc.db.Foreach(bc.group, fn)
}

type BlockchainIterator struct {
	bc   *Blockchain
	next types.HashValue
}

func (bc *Blockchain) Begin() *BlockchainIterator {
	return &BlockchainIterator{bc, bc.GetLastest().Hash()}
}

func (iter *BlockchainIterator) Next() (nextBlock *types.Block) {
	if iter.next == nil {
		return nil
	}
	nextBlock = BytesToBlock(iter.bc.Get(iter.next))
	if nextBlock != nil {
		iter.next = nextBlock.PrevHash
	}
	return nextBlock
}

func CreateBlockchain(minerAddress string) *Blockchain {
	group := global.GetGroup()
	bc := GetBlockchain(group)
	genesis := NewBlock(group, nil, 1, 0, []*types.Transaction{NewCoinbaseTxn(minerAddress)})
	bc.Clear()
	bc.AddBlock(genesis)
	GetUTXOSet(group).Reindex()
	return bc
}

func GetBlockchain(group int) *Blockchain {
	return &Blockchain{global.GetBlocksDB(), group}
}

func (bc *Blockchain) GetGenesis() *types.Block {
	return bc.GetBlockByHeight(0)
}

func (bc *Blockchain) GetLastest() *types.Block {
	block, ok := global.GetBlock(bc.group, "lastest")
	if ok {
		return block
	}
	return BytesToBlock(bc.Get("lastest"))
}

func (bc *Blockchain) GetBlockByHeight(height int32) *types.Block {
	return BytesToBlock(bc.Get(bc.Get(height)))
}

func (bc *Blockchain) SetBlockByHeight(height int, b *types.Block) {
	bc.Set(height, b.Hash())
}

func (bc *Blockchain) SetLastest(bytes []byte) {
	bc.Set("lastest", bytes)
	global.SetBlock(bc.group, "lastest", BytesToBlock(bytes))
}

func (bc *Blockchain) GetHeight() int32 {
	lastest := bc.GetLastest()
	if lastest == nil {
		return -1
	}
	return lastest.Height
}

func (bc *Blockchain) AddBlock(b *types.Block) {
	if b == nil {
		return
	}

	if bc.Get(b.Hash()) != nil {
		return
	}

	bytes := utils.Encode(b)
	if bc.GetHeight() < b.Height {
		bc.SetLastest(bytes)
		bc.Set(b.Hash(), bytes)
		bc.Set(b.Height, b.Hash())
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
		lastDiff := bc.GetBlockByHeight(height - 60)
		thisDiff := bc.GetBlockByHeight(height - 1)
		difficulty *= 59 * 60 * 1e9 / float64(thisDiff.Timestamp-lastDiff.Timestamp)
	}

	newBlock := NewBlock(bc.group, lastest.Hash(), difficulty, height, txns)
	if newBlock == nil {
		return nil
	}
	log.Infof("NewBlock[%d]{%.2f} %s", height, difficulty, lastest.Hash())
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

func (bc *Blockchain) FindTxn(hash types.HashValue) (types.Transaction, error) {
	var block *types.Block
	iter := bc.Begin()

	for block = iter.Next(); block != nil; block = iter.Next() {
		for _, txn := range block.Txns {
			if bytes.Compare(txn.Hash(), hash) == 0 {
				return *txn, nil
			}
		}
	}

	return types.Transaction{}, errors.New(fmt.Sprintf("Transaction is not found, %s", hash))
}

func (bc *Blockchain) SignTransaction(txn *types.Transaction, sk types.PrivateKey) error {
	hashedTxn := make(map[string]types.Transaction)

	for _, vin := range txn.Vin {
		var vinTxn types.Transaction
		vinTxn, err := bc.FindTxn(vin.VoutHash)
		if err != nil {
			return err
		}
		hashedTxn[vinTxn.Hash().String()] = vinTxn
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
		prevTXs[prevTX.Hash().String()] = prevTX
	}

	return txn.Verify(prevTXs)
}
