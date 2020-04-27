package types

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math/big"
	"time"

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
	return fmt.Sprintf("[Left: %v]HashValue: %s", path.Left, path.HashValue)
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

func (b Block) String() (ret string) {
	ret = "\n"
	timeStr := time.Unix(b.Timestamp/1e9, 0).Format("2006/01/02 15:04:05")
	ret += fmt.Sprintf("Group: %d, Height: %d, Timestamp: %s, Target: %f\n",
		b.Group, b.Height, timeStr, b.Target)
	ret += fmt.Sprintf("GroupBase: %d, BatchSize: %d, Nonce: %d\n",
		b.GroupBase, b.BatchSize, b.Nonce)
	ret += fmt.Sprintf("Hash: %s\n", b.Hash())
	ret += fmt.Sprintf("Prev: %s\n", b.PrevHash)
	ret += fmt.Sprintf("MerkleRoot: %s\n", b.MerkleRoot)
	for i, p := range b.BatchMerklePath {
		ret += fmt.Sprintf("BatchMerklePath[%d]: %s\n", i, p)
	}
	for i, txn := range b.Txns {
		ret += fmt.Sprintf("Txn[%d]: %s\n", i, txn)
	}
	return ret
}

func (b Block) Verify() bool {
	div, _ := big.NewFloat(b.Target).Int(nil)
	t := big.NewInt(1)
	target := t.Lsh(t, 256).Div(t, div)

	hashInt := big.NewInt(0)
	hashInt.SetBytes(b.Hash())

	return hashInt.Cmp(target) < 0
}
