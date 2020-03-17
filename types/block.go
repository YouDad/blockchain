package types

import "fmt"

type Block struct {
	Timestamp  int64
	PrevHash   HashValue
	Hash       HashValue
	Difficulty float64
	Nonce      int64
	Height     int32
	MerkleRoot HashValue
	Txns       []*Transaction
}

func (b Block) String() string {
	ret := fmt.Sprintf("\nHeight: %d\n", b.Height)
	ret += fmt.Sprintf("Prev:   %x\n", b.PrevHash)
	ret += fmt.Sprintf("Hash:   %x\n", b.Hash)
	ret += fmt.Sprintf("Txns:\n")
	for i, txn := range b.Txns {
		ret += fmt.Sprintf("\tTxns[%d]:%s", i, txn.String())
	}

	return ret
}
