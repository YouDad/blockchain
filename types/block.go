package types

import (
	"fmt"

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
	return utils.SHA256(b)
}

func (b Block) String() string {
	ret := fmt.Sprintf("\nHeight: %d\n", b.Height)
	ret += fmt.Sprintf("Prev:   %s\n", b.PrevHash)
	ret += fmt.Sprintf("Hash:   %s\n", b.Hash())
	ret += fmt.Sprintf("Txns:\n")
	for i, txn := range b.Txns {
		ret += fmt.Sprintf("\tTxns[%d]:%s", i, txn.String())
	}

	return ret
}
