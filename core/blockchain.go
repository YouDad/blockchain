package core

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/YouDad/blockchain/global"
	"github.com/YouDad/blockchain/global/mempool"
	"github.com/YouDad/blockchain/log"
	"github.com/YouDad/blockchain/types"
	"github.com/YouDad/blockchain/utils"
)

type Blockchain struct {
	db    *global.BlocksDB
	txn   *global.TxnsDB
	group int
}

func (bc *Blockchain) blockClear() {
	bc.db.Clear(bc.group)
}

func (bc *Blockchain) blockGet(key interface{}) (value []byte) {
	return bc.db.Get(bc.group, key)
}

func (bc *Blockchain) blockGetWithoutLog(key interface{}) (value []byte) {
	return bc.db.GetWithoutLog(bc.group, key)
}

func (bc *Blockchain) blockSet(key interface{}, value []byte) {
	bc.db.Set(bc.group, key, value)
}

func (bc *Blockchain) blockDelete(key interface{}) {
	bc.db.Delete(bc.group, key)
}

func (bc *Blockchain) blockForeach(fn func(k, v []byte) bool) {
	bc.db.Foreach(bc.group, fn)
}

func (bc *Blockchain) txnGet(key interface{}) (value []byte) {
	return bc.txn.Get(bc.group, key)
}

func (bc *Blockchain) txnSet(key interface{}, value []byte) {
	bc.txn.Set(bc.group, key, value)
}

func (bc *Blockchain) txnDelete(key interface{}) {
	bc.txn.Delete(bc.group, key)
}

func (bc *Blockchain) txnForeach(fn func(k, v []byte) bool) {
	bc.txn.Foreach(bc.group, fn)
}

func (bc *Blockchain) TxnReindex() {
	bc.txn.Clear(bc.group)
	iter := bc.Begin()
	for block := iter.Next(); block != nil; block = iter.Next() {
		for _, txn := range block.Txns {
			bc.txnSet(txn.Hash(), utils.Encode(txn))
		}
	}
}

type BlockchainIterator struct {
	bc   *Blockchain
	next types.HashValue
}

func (bc *Blockchain) Begin() *BlockchainIterator {
	log.Debugln("Blockchain.Begin")
	lastest := BytesToBlock(bc.blockGet("lastest"))
	return &BlockchainIterator{bc, lastest.Hash()}
}

func (iter *BlockchainIterator) Next() (nextBlock *types.Block) {
	if iter.next == nil {
		return nil
	}
	nextBlock = BytesToBlock(iter.bc.blockGetWithoutLog(iter.next))
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
	bc.blockClear()
	bc.AddBlock(block)
	GetUTXOSet(group).Reindex()
	bc.TxnReindex()
	return nil
}

func GetBlockchain(group int) *Blockchain {
	return &Blockchain{
		db:    global.GetBlocksDB(),
		txn:   global.GetTxnsDB(),
		group: group % global.MaxGroupNum,
	}
}

func (bc *Blockchain) GetGenesis() *types.Block {
	hash := bc.blockGet(0)
	if hash == nil {
		hash = bc.blockGet(0)
	}

	block := bc.blockGet(hash)
	if block == nil {
		block = bc.blockGet(hash)
	}

	return BytesToBlock(block)
}

func (bc *Blockchain) GetLastest() *types.Block {
	return BytesToBlock(bc.blockGet("lastest"))
}

func (bc *Blockchain) GetBlockByHash(hash types.HashValue) *types.Block {
	return BytesToBlock(bc.blockGet(hash))
}

func (bc *Blockchain) GetBlockByHeight(height int32) *types.Block {
	hash := bc.blockGet(height)
	if hash == nil {
		hash = bc.blockGet(height)
	}

	block := bc.blockGet(hash)
	if block == nil {
		block = bc.blockGet(hash)
	}

	return BytesToBlock(block)
}

func (bc *Blockchain) SetBlockByHeight(height int, b *types.Block) {
	bc.blockSet(height, b.Hash())
}

func (bc *Blockchain) SetLastest(block *types.Block) {
	bc.blockSet("lastest", utils.Encode(block))
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
		defer mutexHeight.Unlock()
		if lastest == nil {
			cacheHeight[bc.group] = -1
		} else {
			cacheHeight[bc.group] = lastest.Height
		}
	})
	mutexOnceGetHeight.Unlock()
	mutexHeight.Lock()
	defer mutexHeight.Unlock()
	return cacheHeight[bc.group]
}

func (bc *Blockchain) AddBlock(b *types.Block) {
	// 去重
	if b == nil {
		return
	}

	// 符合高度
	if bc.GetHeight()+1 == b.Height {
		bytes := utils.Encode(b)
		bc.SetLastest(b)
		bc.blockSet(b.Hash(), bytes)
		bc.blockSet(b.Height, b.Hash())
		GetBlockhead(bc.group).AddBlockhead(b)
		for _, txn := range b.Txns {
			bc.txnSet(txn.Hash(), utils.Encode(txn))
		}

		m := mempool.GetMempool(b.Group)
		for _, txn := range b.Txns {
			m.Delete(txn.Hash())
		}

		// TODO: 还可以优化写法
		var verify func(txn *types.Transaction) bool
		verify = func(txn *types.Transaction) bool {
			for _, vin := range txn.Vin {
				txnHash := vin.VoutHash
				_, err := bc.FindTxn(txnHash)
				if err != nil {
					t, err := m.GetTxn(txnHash)
					if err != nil || !verify(t) {
						return false
					}
				}
			}
			return true
		}

		txns := m.GetTxns()
		for _, txn := range txns {
			if !verify(txn) {
				log.Warnln("Delete Txn !!!!!!")
				m.Delete(txn.Hash())
			}
		}
		m.Release()
	}
}

func (bc *Blockchain) DeleteBlock(b *types.Block) {
	bc.blockDelete(b.Hash())
	bc.blockDelete(b.Height)
	GetBlockhead(bc.group).Delete(b.Height)
	for _, txn := range b.Txns {
		bc.txnDelete(txn.Hash())
	}
}

func MineBlocksForCreate(txn *types.Transaction, groupBase int) (*types.Block, error) {
	txns := []*types.Transaction{txn}
	blocks := []*types.Block{{
		BlockHeader: types.BlockHeader{
			Group:      groupBase,
			Height:     0,
			PrevHash:   nil,
			Timestamp:  time.Now().UnixNano(),
			MerkleRoot: NewTxnMerkleTree(txns).RootNode.Data,
			Target:     1.0,
		},
		ChukonuHeader: types.ChukonuHeader{
			GroupBase: groupBase,
			BatchSize: 1,
		},
		Txns: txns,
	}}

	pow := NewPOW(blocks)
	err, nonce, _ := pow.Run()
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
			return nil, errors.New(fmt.Sprintf(
				"MineBlocks failed, because don't have "+
					"blockchain[%d].lastest", (groupBase+i)%global.MaxGroupNum))
		}
		var target float64 = lastest.Target
		var height int32 = lastest.Height
		var timestamp int64 = lastest.Timestamp

		// 1. 更新难度
		if height%30 == 0 {
			prevRecalcBlock := bc.GetBlockByHeight(height - 29)
			if prevRecalcBlock != nil {
				target *= 29 * 30 * 1e9 / float64(timestamp-prevRecalcBlock.Timestamp)
			}
		}

		// 2. 构造block
		blocks = append(blocks, &types.Block{
			BlockHeader: types.BlockHeader{
				Group:      (groupBase + i) % global.MaxGroupNum,
				Height:     height + 1,
				PrevHash:   lastest.Hash(),
				Timestamp:  time.Now().UnixNano(),
				MerkleRoot: NewTxnMerkleTree(txns[i]).RootNode.Data,
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
	err, nonce, batchMerkleTree := pow.Run()
	if err != nil {
		return nil, err
	}

	// 3. 过滤有效区块
	var foundBlocks []*types.Block
	for _, block := range blocks {
		block.Nonce = nonce
		block.BatchMerklePath = batchMerkleTree.FindPath(
			(block.Group + global.MaxGroupNum - block.GroupBase) % global.MaxGroupNum)
		if block.Verify() {
			foundBlocks = append(foundBlocks, block)
		}
	}

	// 4. 返回挖到的区块
	for _, block := range foundBlocks {
		log.Infoln("NewBlock", block.Hash(), block)
	}
	return foundBlocks, nil
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
	b := bc.txnGet(hash)
	if len(b) == 0 {
		return nil, errors.New(fmt.Sprintf("Transaction is not found, %s", hash))
	}
	return BytesToTransaction(b), nil
}

func (bc *Blockchain) SignTransaction(txn *types.Transaction, sk types.PrivateKey) error {
	hashedTxn := make(map[string]types.Transaction)

	for _, vin := range txn.Vin {
		vinTxn, err := bc.FindTxn(vin.VoutHash)
		if err == nil {
			hashedTxn[vinTxn.Hash().String()] = *vinTxn
			continue
		}

		vinTxn, err = mempool.GetTxn(bc.group, vin.VoutHash)
		if err != nil {
			return errors.New(fmt.Sprintf("Transaction is not found, %s", vin.VoutHash))
		}
		hashedTxn[vinTxn.Hash().String()] = *vinTxn
	}

	return txn.Sign(sk, hashedTxn)
}

// 验证交易是否有效
func (bc *Blockchain) VerifyTransaction(txn types.Transaction) error {
	if txn.IsCoinbase() {
		return nil
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
		prevTxn, err = mempool.GetTxn(bc.group, vin.VoutHash)
		if err != nil {
			return errors.New("mempool not found " + vin.VoutHash.String())
		}
		prevTxns[prevTxn.Hash().String()] = *prevTxn
	}

	if txn.Verify(prevTxns) {
		return nil
	} else {
		return errors.New("verify false")
	}
}
