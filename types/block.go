package types

import (
	"fmt"

	"github.com/YouDad/blockchain/utils"
)

type Block struct {
	Timestamp  int64
	PrevHash   HashValue
	Difficulty float64
	Nonce      int64
	MerkleRoot HashValue
	Height     int32
	Txns       []*Transaction
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
