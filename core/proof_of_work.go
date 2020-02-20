package core

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"log"
	"math"
	"math/big"

	"github.com/YouDad/blockchain/utils"
)

const targetBits = 12

type ProofOfWork struct {
	block  *Block
	target *big.Int
}

func NewProofOfWork(b *Block) *ProofOfWork {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-targetBits))
	return &ProofOfWork{b, target}
}

func (pow *ProofOfWork) prepareData(nonce int64) []byte {
	var txs [][]byte

	for _, tx := range pow.block.App.Transactions {
		txs = append(txs, tx.Serialize())
	}
	mTree := NewMerkleTree(txs)

	return bytes.Join(
		[][]byte{
			pow.block.PrevBlockHash,
			mTree.RootNode.Data,
			utils.IntToHex(pow.block.Timestamp),
			utils.IntToHex(int64(targetBits)),
			utils.IntToHex(int64(nonce)),
		},
		[]byte{},
	)
}

func (pow *ProofOfWork) Run() (int64, []byte) {
	var hashInt big.Int
	var hash [32]byte
	var nonce int64 = 0

	fmt.Printf("Mining a new block %s\n", pow.block.App.ToString())

	for nonce < math.MaxInt64 {
		data := pow.prepareData(nonce)
		hash = sha256.Sum256(data)
		if nonce%(1<<16) == 0 {
			log.Printf("Dig into mine [%d] %x\n", nonce, hash)
		}

		hashInt.SetBytes(hash[:])
		if hashInt.Cmp(pow.target) == -1 {
			break
		} else {
			nonce++
		}
	}
	fmt.Print("\n\n")
	return nonce, hash[:]
}

func (pow *ProofOfWork) Validate() bool {
	var hashInt big.Int

	data := pow.prepareData(pow.block.Nonce)
	hash := sha256.Sum256(data)
	hashInt.SetBytes(hash[:])

	isValid := hashInt.Cmp(pow.target) == -1

	return isValid
}
