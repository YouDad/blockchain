package core

import (
	"bytes"
	"crypto/sha256"
	"math"
	"math/big"
	"math/rand"
	"time"

	"github.com/YouDad/blockchain/global"
	"github.com/YouDad/blockchain/log"
	"github.com/YouDad/blockchain/types"
	"github.com/YouDad/blockchain/utils"
)

var hashSpeed uint = 100

func Register(speed uint) {
	hashSpeed = speed
}

type ProofOfWork struct {
	block  *types.Block
	target *big.Int
}

func NewPOW(b *types.Block) *ProofOfWork {
	max := big.NewInt(1)
	max.Lsh(max, 256)
	target := big.NewInt(0)
	big.NewFloat(b.Target).Int(target)
	max.Div(max, target)
	log.Debugln(max, target)
	return &ProofOfWork{b, max}
}

func (pow *ProofOfWork) prepareData(nonce int64) []byte {
	return bytes.Join(
		[][]byte{
			pow.block.PrevHash[:],
			pow.block.MerkleRoot,
			utils.IntToBytes(pow.block.Timestamp),
			utils.FloatToBytes(pow.block.Target),
			utils.IntToBytes(nonce),
		},
		[]byte{},
	)
}

func (pow *ProofOfWork) Run(group int) (int64, []byte) {
	var hash []byte
	var ok bool
	rand.Seed(time.Now().UnixNano())
	var nonce int64 = rand.Int63()
	start := nonce

	for nonce < math.MaxInt64 {
		if GetBlockchain(group).GetHeight() != pow.block.Height-1 {
			return 0, nil
		}
		ok, hash = pow.Validate(nonce)
		if ok {
			break
		} else {
			nonce++
		}

		if (nonce-start)%int64(hashSpeed*1e4) == 0 {
			if hashSpeed < 100 {
				time.Sleep(time.Duration(1e7 * (100 - hashSpeed)))
			}
		}
	}
	global.SyncMutex.Lock()
	return nonce, hash
}

func (pow *ProofOfWork) Validate(nonce int64) (bool, types.HashValue) {
	var hashInt big.Int

	data := pow.prepareData(nonce)
	hash := sha256.Sum256(data)
	hashInt.SetBytes(hash[:])

	isValid := hashInt.Cmp(pow.target) == -1

	if nonce%(1<<20) == 0 {
		log.Debugf("Dig into mine [%d] %s\n", nonce, hash)
	}

	return isValid, hash[:]
}
