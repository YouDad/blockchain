package core

import (
	"encoding/hex"

	"github.com/YouDad/blockchain/log"
	"github.com/YouDad/blockchain/store"
	"github.com/YouDad/blockchain/types"
	"github.com/YouDad/blockchain/utils"
)

type Blockchain struct {
	store.Database
}

func CreateBlockchain(minerAddress string) Blockchain {
	bc := Blockchain{store.CreateDatabase()}
	bc.SetTable("Blocks").Clear()
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
	log.NotImplement()
	return nil
}

func (bc *Blockchain) GetHeight() int {
	log.NotImplement()
	return 0
}

func (bc *Blockchain) AddBlock(b *Block) {
	log.NotImplement()
}

func (bc *Blockchain) MineBlock(txns []*Transaction) *Block {
	log.NotImplement()
	return nil
}

type BlockchainIterator struct {
	*Blockchain
	next types.HashValue
}

func (bc *Blockchain) Begin() *BlockchainIterator {
	lastestBlock := BytesToBlock(bc.SetTable("Blocks").Get("lastest"))
	return &BlockchainIterator{bc, lastestBlock.Hash()}
}

func (iter *BlockchainIterator) Next() (nextBlock *Block) {
	iter.SetTable("Blocks")
	if iter.next == nil {
		return nil
	}
	nextBlock = BytesToBlock(iter.Get(iter.next))
	iter.next = nextBlock.PrevHash
	return nextBlock
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
