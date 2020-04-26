package core

import (
	"bytes"
	"errors"
	"fmt"
	"sync"
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

func (bc *Blockchain) GetWithoutLog(key interface{}) (value []byte) {
	return bc.db.GetWithoutLog(bc.group, key)
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
	log.Debugln("Blockchain.Begin")
	return &BlockchainIterator{bc, bc.GetLastest().Hash()}
}

func (iter *BlockchainIterator) Next() (nextBlock *types.Block) {
	if iter.next == nil {
		return nil
	}
	nextBlock = BytesToBlock(iter.bc.GetWithoutLog(iter.next))
	if nextBlock != nil {
		iter.next = nextBlock.PrevHash
	}
	return nextBlock
}

func CreateBlockchain(minerAddress string) error {
	group := global.GetGroup()
	bc := GetBlockchain(group)
	block, err := MineBlocksForCreate(NewCoinbaseTxn(minerAddress), group)
	if err != nil {
		return err
	}
	bc.Clear()
	bc.AddBlock(block)
	GetUTXOSet(group).Reindex()
	return nil
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
	block := BytesToBlock(bytes)
	global.SetBlock(bc.group, "lastest", block)
	mutexHeight.Lock()
	cacheHeight[bc.group] = block.Height
	mutexHeight.Unlock()
}

var (
	onceGetHeight      = make(map[int]*sync.Once)
	cacheHeight        = make(map[int]int32)
	mutexHeight        sync.Mutex
	mutexOnceGetHeight sync.Mutex
)

func (bc *Blockchain) GetHeight() int32 {
	mutexOnceGetHeight.Lock()
	_, ok := onceGetHeight[bc.group]
	if !ok {
		onceGetHeight[bc.group] = &sync.Once{}
	}

	onceGetHeight[bc.group].Do(func() {
		lastest := bc.GetLastest()
		mutexHeight.Lock()
		if lastest == nil {
			cacheHeight[bc.group] = -1
		} else {
			cacheHeight[bc.group] = lastest.Height
		}
		mutexHeight.Unlock()
	})
	mutexOnceGetHeight.Unlock()
	return cacheHeight[bc.group]
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

		for _, txn := range b.Txns {
			global.GetMempool(b.Group).Delete(txn.Hash())
		}
	}
}

func MineBlocksForCreate(txns *types.Transaction, groupBase int) (*types.Block, error) {
	blocks := []*types.Block{{
		BlockHeader: types.BlockHeader{
			Group:      groupBase,
			Height:     0,
			PrevHash:   nil,
			Timestamp:  time.Now().UnixNano(),
			MerkleRoot: NewMerkleTree([][]byte{utils.Encode(txns)}).RootNode.Data,
			Target:     1.0,
		},
		ChukonuHeader: types.ChukonuHeader{
			GroupBase: groupBase,
			BatchSize: 1,
		},
		Txns: []*types.Transaction{txns},
	}}

	pow := NewPOW(blocks)
	err, _, nonce, _ := pow.Run()
	if err != nil {
		return nil, err
	}

	blocks[0].Nonce = nonce
	log.Infoln("CreateBlockchain", blocks[0])
	return blocks[0], nil
}

func MineBlocks(txns [][]*types.Transaction, groupBase, batchSize int) ([]*types.Block, error) {
	// 1. 构造blocks
	var blocks []*types.Block
	for i := 0; i < batchSize; i++ {
		bc := GetBlockchain(groupBase + i)
		lastest := bc.GetLastest()
		if lastest == nil {
			return nil, errors.New(fmt.Sprintf("MineBlocks failed, because don't have blockchain[%d].lastest", groupBase+i))
		}
		var target float64 = lastest.Target
		var height int32 = lastest.Height
		var prevHash types.HashValue = lastest.Hash()
		var timestamp int64 = lastest.Timestamp

		// 1. 计算MerkleRoot用字节数组
		var txnsBytes [][]byte
		for _, txn := range txns[i] {
			txnsBytes = append(txnsBytes, utils.Encode(txn))
		}

		// 2. 更新难度
		if height%30 == 0 {
			prevRecalcBlock := bc.GetBlockByHeight(height - 29)
			if prevRecalcBlock != nil {
				target *= 29 * 30 * 1e9 / float64(timestamp-prevRecalcBlock.Timestamp)
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
		return nil, err
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
		log.Infoln("NewBlock", block)
	}
	return blocks, nil
}

func (bc *Blockchain) FindUTXO() map[string][]types.TxnOutput {
	utxos := make(map[string][]types.TxnOutput)
	stxos := make(map[string]map[int]bool)

	// 遍历区块链
	iter := bc.Begin()
	for {
		block := iter.Next()
		if block == nil {
			break
		}

		// 遍历区块的所有交易
		for _, txn := range block.Txns {
			hash := txn.Hash().String()

			// 遍历所有输出
			for i, out := range txn.Vout {
				_, ok := stxos[hash]
				if !ok {
					utxos[hash] = append(utxos[hash], out)
					continue
				}
				_, ok = stxos[hash][i]
				if !ok {
					utxos[hash] = append(utxos[hash], out)
				}
			}

			// 遍历所有输入
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

func (bc *Blockchain) FindTxn(hash types.HashValue) (*types.Transaction, error) {
	var block *types.Block
	iter := bc.Begin()

	// 遍历所有区块的所有交易，找到哈希值一致的交易
	for block = iter.Next(); block != nil; block = iter.Next() {
		for _, txn := range block.Txns {
			if bytes.Compare(txn.Hash(), hash) == 0 {
				return txn, nil
			}
		}
	}

	return nil, errors.New(fmt.Sprintf("Transaction is not found, %s", hash))
}

func (bc *Blockchain) SignTransaction(txn *types.Transaction, sk types.PrivateKey) error {
	hashedTxn := make(map[string]types.Transaction)

	for _, vin := range txn.Vin {
		vinTxn, err := bc.FindTxn(vin.VoutHash)
		if err == nil {
			hashedTxn[vinTxn.Hash().String()] = *vinTxn
			continue
		}

		vinTxn, err = global.GetMempool(bc.group).GetTxn(vin.VoutHash)
		if err != nil {
			return errors.New(fmt.Sprintf("Transaction is not found, %s", vin.VoutHash))
		}
		hashedTxn[vinTxn.Hash().String()] = *vinTxn
	}

	return txn.Sign(sk, hashedTxn)
}

// 验证交易是否有效
func (bc *Blockchain) VerifyTransaction(txn types.Transaction) bool {
	if txn.IsCoinbase() {
		return true
	}

	// map[前置交易哈希]前置交易
	prevTxns := make(map[string]types.Transaction)

	// 遍历交易的输入
	for _, vin := range txn.Vin {
		// 在区块链中用引用交易哈希找到该输入引用的交易
		prevTxn, err := bc.FindTxn(vin.VoutHash)
		if err == nil {
			prevTxns[prevTxn.Hash().String()] = *prevTxn
			continue
		}

		// 在未打包交易池中用引用交易哈希找到该输入引用的交易
		prevTxn, err = global.GetMempool(bc.group).GetTxn(vin.VoutHash)
		if err != nil {
			log.Traceln("Not Found", vin.VoutHash)
			return false
		}
		prevTxns[prevTxn.Hash().String()] = *prevTxn
	}

	return txn.Verify(prevTxns)
}
