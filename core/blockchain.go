package core

import (
	"bytes"
	"errors"
	"fmt"
	"time"

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
	txns := [][]*types.Transaction{[]*types.Transaction{NewCoinbaseTxn(minerAddress)}}
	blocks := MineBlocks(txns, group, 1)
	bc.Clear()
	bc.AddBlock(blocks[0])
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
	// 去重
	if b == nil || bc.Get(b.Hash()) != nil {
		return
	}

	// 符合高度
	if bc.GetHeight()+1 == b.Height {
		bytes := utils.Encode(b)
		bc.SetLastest(bytes)
		bc.Set(b.Hash(), bytes)
		bc.Set(b.Height, b.Hash())
		GetBlockhead(bc.group).AddBlockhead(b)

		mempool := global.GetMempool(b.Group)
		for _, txn := range b.Txns {
			delete(mempool, txn.Hash().Key())
		}
	}
}

func MineBlocks(txns [][]*types.Transaction, groupBase, batchSize int) []*types.Block {
	// 1. 构造blocks
	var blocks []*types.Block
	for i := 0; i < batchSize; i++ {
		bc := GetBlockchain(groupBase + i)
		lastest := bc.GetLastest()
		var target float64 = 1.0
		var height int32 = -1
		var prevHash types.HashValue = nil
		var timestamp int64 = time.Now().UnixNano()
		if lastest != nil {
			target = lastest.Target
			height = lastest.Height
			prevHash = lastest.Hash()
			timestamp = lastest.Timestamp
		}

		// 1. 计算MerkleRoot用字节数组
		var txnsBytes [][]byte
		for _, txn := range txns[i] {
			txnsBytes = append(txnsBytes, utils.Encode(txn))
		}

		// 2. 更新难度
		if height%60 == 0 {
			prevRecalcBlock := bc.GetBlockByHeight(height - 59)
			if prevRecalcBlock != nil {
				target *= 59 * 60 * 1e9 / float64(timestamp-prevRecalcBlock.Timestamp)
			}
		}

		// 3. 构造block
		blocks = append(blocks, &types.Block{
			BlockHeader: types.BlockHeader{
				Group:      groupBase + i,
				Height:     height + 1,
				PrevHash:   prevHash,
				Timestamp:  time.Now().UnixNano(),
				MerkleRoot: NewMerkleTree(txnsBytes).RootNode.Data,
				Target:     target,
			},
			ChukonuHeader: types.ChukonuHeader{
				GroupBase: groupBase,
				BatchSize: batchSize,
			},
			Txns: txns[i],
		})
	}

	// 2. 计算Nonce
	pow := NewPOW(blocks)
	err, target, nonce, batchMerkleTree := pow.Run()
	if err != nil {
		return nil
	}

	// 3. 过滤有效区块
	var foundBlocks []*types.Block
	for _, block := range blocks {
		if target.Cmp(GetTarget(block.Target)) == -1 {
			foundBlocks = append(foundBlocks, block)
		}
	}

	// 4. 设置区块的一些字段
	for _, block := range foundBlocks {
		block.Nonce = nonce
		block.BatchMerklePath = batchMerkleTree.FindPath(block.Group - block.GroupBase)
	}

	// 5. 返回挖到的区块
	for _, block := range blocks {
		log.Infof("NewBlock[%d]{%.2f} %s", block.Height, block.Target, block.PrevHash)
	}
	return blocks
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
