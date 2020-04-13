package api

import (
	"errors"

	"github.com/YouDad/blockchain/apicb"
	"github.com/YouDad/blockchain/core"
	"github.com/YouDad/blockchain/log"
	"github.com/YouDad/blockchain/rpc"
	"github.com/YouDad/blockchain/types"
)

var (
	NULL                   = true
	RootHashDifferentError = errors.New("RootHash is different.")
	VersionDifferentError  = errors.New("Version is different.")
)

func SendVersion(nowHeight int, genesisHash types.HashValue) (height int, err error) {
	log.Errln("NotImplement")
	return 0, nil
}

func GetGenesis() (*core.Block, error) {
	log.Errln("NotImplement")
	var genesisBlock apicb.GetGenesisReply
	err := rpc.Call("DBApi.GetGenesis", &NULL, &genesisBlock)
	if err != nil {
		return nil, err
	}
	return nil, err
}

func GetBlocks(start, end int) []*core.Block {
	log.Errln("NotImplement")
	return nil
}

func GetBalance(address string) (int, error) {
	log.Errln("NotImplement")
	return 0, nil
}

func GetVersion() (types.Version, error) {
	log.Errln("NotImplement")
	return types.Version{}, nil
}

func SendTransaction(txn *core.Transaction) {
	log.Errln("NotImplement")

}
