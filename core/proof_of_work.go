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
	"github.com/YouDad/blockchain/types"
	"github.com/YouDad/blockchain/utils"
)

const targetBits int64 = 16

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
	var hash []byte
	var ok bool
	rand.Seed(time.Now().UnixNano())
	var nonce int64 = rand.Int63()

	for nonce < math.MaxInt64 {
		ok, hash = pow.Validate(nonce)
		if ok {
			break
		} else {
			nonce++
		}
	}
	fmt.Print("\n\n")
	return nonce, hash
}

func (pow *ProofOfWork) Validate(nonce int64) (bool, types.HashValue) {
	var hashInt big.Int

	data := pow.prepareData(nonce)
	hash := sha256.Sum256(data)
	hashInt.SetBytes(hash[:])

	isValid := hashInt.Cmp(pow.target) == -1

	if nonce%(1<<16) == 0 {
		log.Infof("Dig into mine [%d] %x\n", nonce, hash)
	}

	return isValid, hash[:]
}
