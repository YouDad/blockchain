package core

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math"
	"math/big"
	"math/rand"
	"time"

	"github.com/YouDad/blockchain/log"
	"github.com/YouDad/blockchain/utils"
)

const targetBits int64 = 12

type ProofOfWork struct {
	block  *Block
	target *big.Int
}

func NewPOW(b *Block) *ProofOfWork {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-targetBits))
	return &ProofOfWork{b, target}
}

func (pow *ProofOfWork) prepareData(nonce int64) []byte {
	return bytes.Join(
		[][]byte{
			pow.block.PrevHash[:],
			pow.block.MerkleRoot,
			utils.IntToBytes(pow.block.Timestamp),
			utils.IntToBytes(targetBits),
			utils.IntToBytes(nonce),
		},
		[]byte{},
	)
}

func (pow *ProofOfWork) Run() (int64, []byte) {
	var hashInt big.Int
	var hash [32]byte
	rand.Seed(time.Now().UnixNano())
	var nonce int64 = rand.Int63()

	log.Infof("Mining a new block %s\n", pow.block.String())

	for nonce < math.MaxInt64 {
		data := pow.prepareData(nonce)
		hash = sha256.Sum256(data)
		if nonce%(1<<16) == 0 {
			log.Infof("Dig into mine [%d] %x\n", nonce, hash)
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
