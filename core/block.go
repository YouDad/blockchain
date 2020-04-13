package core

import (
	"github.com/YouDad/blockchain/log"
	"github.com/YouDad/blockchain/types"
	"github.com/YouDad/blockchain/utils"
)

type BlockHead struct {
	Timestamp  int64
	PrevHash   types.HashValue
	Nonce      int64
	Height     int
	MerkleRoot types.HashValue
}

func (bh *BlockHead) GobEncode() (bytes []byte, err error) {
	log.Errln("NotImplement")
	return nil, nil
}

func (bh *BlockHead) GobDecode(bytes []byte) error {
	log.Errln("NotImplement")
	return nil
}

type Block struct {
	BlockHead
	Txns []*Transaction
}

func (bh *BlockHead) Hash() types.HashValue {
	return utils.SHA256(bh)
}
