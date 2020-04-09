package types

import (
	"bytes"
	"crypto/sha256"
	"fmt"
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

type ChukonuHeader struct {
	GroupBase       int
	BatchSize       int
	BatchMerklePath []HashValue
	Nonce           int64
}

type Block struct {
	BlockHeader
	ChukonuHeader
	Txns []*Transaction
}

func (b Block) Hash() HashValue {
	hash := b.BlockHeader.Hash()
	for _, hashValue := range b.BatchMerklePath {
		sha := sha256.Sum256(append(hash, hashValue...))
		hash = sha[:]
	}
	sha := sha256.Sum256(bytes.Join(
		[][]byte{
			hash,
			utils.BaseTypeToBytes(b.GroupBase),
			utils.BaseTypeToBytes(b.BatchSize),
			utils.BaseTypeToBytes(b.Nonce),
		}, []byte{},
	))
	return sha[:]
}

func (b Block) String() string {
	ret := fmt.Sprintf("[%d]", b.Height)
	ret += fmt.Sprintf("%s<-", b.PrevHash[:3])
	ret += fmt.Sprintf("%s", b.Hash()[:3])
	return ret
}

func (b Block) Verify() bool {
	hashInt := big.NewInt(0)
	target := big.NewInt(1)
	hashInt.SetBytes(b.Hash())
	div, _ := big.NewFloat(b.Target).Int(nil)
	return hashInt.Cmp(target.Lsh(target, 256).Div(target, div)) < 0
}
