package api

import (
	"errors"

	"github.com/YouDad/blockchain/core"
	"github.com/YouDad/blockchain/log"
	"github.com/YouDad/blockchain/network"
	"github.com/YouDad/blockchain/types"
)

var (
	RootHashDifferentError = errors.New("RootHash is different.")
	VersionDifferentError  = errors.New("Version is different.")
)

func SendVersion(nowHeight int, genesisHash types.HashValue) (height int, err error) {
	log.NotImplement()
	return 0, nil
}

func GetGenesis() (*core.Block, error) {
	log.NotImplement()
	var genesisBlock GetGenesisReply
	err := network.Call("DB.GetGenesis", &NULL, &genesisBlock)
	if err != nil {
		return nil, err
	}
	return nil, err
}

func GetBlocks(start, end int) []*core.Block {
	log.NotImplement()
	return nil
}

func GetBalance(address string) (balance GetBalanceReply, err error) {
	err = network.CallMySelf("DB.GetBalance", &address, &balance)
	return balance, err
}

func SendTransaction(txn *core.Transaction) {
	log.NotImplement()

}
