package types

import (
	"bytes"
	"crypto/sha256"
	"math/big"

	"github.com/YouDad/blockchain/utils"
)

type BlockHeader struct {
	Group      int
	Height     int32
	PrevHash   HashValue
	Timestamp  int64
	MerkleRoot HashValue
	Target     float64
}

func (bh BlockHeader) Hash() HashValue {
	return utils.SHA256(bh)
}

type MerklePath struct {
	HashValue HashValue
	Left      bool
}

func (path MerklePath) String() string {
	return string(utils.Encode(path))
}

type ChukonuHeader struct {
	GroupBase       int
	BatchSize       int
	BatchMerklePath []MerklePath
	Nonce           int64
}

type Block struct {
	BlockHeader
	ChukonuHeader
	Txns []*Transaction
}

func (b Block) Hash() HashValue {
	var hash HashValue

	sha := sha256.Sum256(b.BlockHeader.Hash())
	hash = sha[:]
	for _, path := range b.BatchMerklePath {
		if path.Left {
			sha = sha256.Sum256(append(path.HashValue, hash...))
		} else {
			sha = sha256.Sum256(append(hash, path.HashValue...))
		}
		hash = sha[:]
	}

	sha = sha256.Sum256(bytes.Join(
		[][]byte{
			sha[:],
			utils.BaseTypeToBytes(b.GroupBase),
			utils.BaseTypeToBytes(b.BatchSize),
			utils.BaseTypeToBytes(b.Nonce),
		}, []byte{},
	))
	return sha[:]
}

func (b Block) String() string {
	return string(utils.Encode(b))
}

func (b Block) Verify() bool {
	div, _ := big.NewFloat(b.Target).Int(nil)
	t := big.NewInt(1)
	target := t.Lsh(t, 256).Div(t, div)

	hashInt := big.NewInt(0)
	hashInt.SetBytes(b.Hash())

	return hashInt.Cmp(target) < 0
}
