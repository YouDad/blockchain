package core

import (
	"bytes"
	"crypto/sha256"
	"errors"
	"math"
	"math/big"
	"math/rand"
	"time"

	"github.com/YouDad/blockchain/log"
	"github.com/YouDad/blockchain/types"
	"github.com/YouDad/blockchain/utils"
)

var (
	hashSpeed           uint = 100
	nonceStart          int64
	ErrBlockchainChange = errors.New("blockchain is changed")
)

func Register(speed uint) {
	hashSpeed = speed
}

type PoweredStruct struct {
	batchMerkleRoot types.HashValue
	GroupBase       int
	BatchSize       int
	Nonce           int64
}

type ProofOfWork struct {
	blocks          []*types.Block
	target          *big.Int
	batchMerkleTree *MerkleTree
	poweredStruct   PoweredStruct
}

func NewPOW(blocks []*types.Block) *ProofOfWork {
	pow := &ProofOfWork{
		blocks: blocks,
		poweredStruct: PoweredStruct{
			GroupBase: blocks[0].Group,
			BatchSize: len(blocks),
		},
	}

	// 1. Calc target
	minTarget := math.MaxFloat64
	for _, block := range blocks {
		if minTarget > block.Target {
			minTarget = block.Target
		}
	}
	log.Debugln(minTarget)
	pow.target = GetTarget(minTarget)

	// 2. Random for nonce
	rand.Seed(time.Now().UnixNano())
	pow.poweredStruct.Nonce = rand.Int63()
	nonceStart = pow.poweredStruct.Nonce

	// 3. get merkle tree
	var nodeSeq [][]byte
	for _, block := range blocks {
		nodeSeq = append(nodeSeq, block.Hash())
	}
	pow.batchMerkleTree = NewMerkleTree(nodeSeq)
	pow.poweredStruct.batchMerkleRoot = pow.batchMerkleTree.RootNode.Data

	return pow
}

func GetTarget(target float64) *big.Int {
	div, _ := big.NewFloat(target).Int(nil)
	t := big.NewInt(1)
	return t.Lsh(t, 256).Div(t, div)
}

func (pow *ProofOfWork) Run() (error, *big.Int, int64, *MerkleTree) {
	nonce := pow.poweredStruct.Nonce
	var target *big.Int

	for nonce < math.MaxInt64 {
		for _, block := range pow.blocks {
			if GetBlockchain(block.Group).GetHeight() != block.Height-1 {
				return ErrBlockchainChange, nil, 0, nil
			}
		}

		ok, t := pow.Validate(nonce)
		if ok {
			target = t
			break
		} else {
			nonce++
		}

		if (nonce-nonceStart)%int64(hashSpeed*1e4) == 0 {
			if hashSpeed < 100 {
				time.Sleep(time.Duration(1e7 * (100 - hashSpeed)))
			}
		}
	}
	return nil, target, nonce, pow.batchMerkleTree
}

func (pow *ProofOfWork) Validate(nonce int64) (bool, *big.Int) {
	hashInt := big.NewInt(0)

	pow.poweredStruct.Nonce = nonce
	data := pow.prepareData()
	hash := sha256.Sum256(data)
	hashInt.SetBytes(hash[:])

	if (nonce-nonceStart)%(1<<20) == 0 {
		log.Debugf("Dig into mine [%d] %s\n", nonce-nonceStart, hash)
	}
	return hashInt.Cmp(pow.target) == -1, hashInt
}

func (pow *ProofOfWork) prepareData() []byte {
	return bytes.Join(
		[][]byte{
			pow.poweredStruct.batchMerkleRoot,
			utils.BaseTypeToBytes(pow.poweredStruct.GroupBase),
			utils.BaseTypeToBytes(pow.poweredStruct.BatchSize),
			utils.BaseTypeToBytes(pow.poweredStruct.Nonce),
		}, []byte{},
	)
}
