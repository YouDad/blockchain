package api

import (
	"errors"
	"net/rpc"

	"github.com/YouDad/blockchain/conf"
	"github.com/YouDad/blockchain/core"
	"github.com/YouDad/blockchain/log"
	"github.com/YouDad/blockchain/network"
	"github.com/YouDad/blockchain/store"
	"github.com/YouDad/blockchain/utils"
	"github.com/YouDad/blockchain/wallet"
)

type DB struct{}

func init() {
	log.Err(rpc.Register(new(DB)))
}

type GetGenesisArgs = conf.NIL
type GetGenesisReply = core.Block

var (
	ErrNull = errors.New("null")
)

func GetGenesis() (*core.Block, error) {
	var genesis GetGenesisReply
	err := network.Call("DB.GetGenesis", &conf.NULL, &genesis)
	if err != nil {
		return nil, err
	}
	return &genesis, err
}

func (db *DB) GetGenesis(args *GetGenesisArgs, reply *GetGenesisReply) error {
	log.Infoln("Response GetGenesis")
	if store.IsDatabaseExists() {
		genesis := core.GetBlockchain().GetGenesis()
		*reply = *genesis
		return nil
	}
	return ErrNull
}

type GetBalanceArgs = string
type GetBalanceReply = int64

func GetBalance(address string) (balance GetBalanceReply, err error) {
	err = network.CallMySelf("DB.GetBalance", &address, &balance)
	return balance, err
}

func (db *DB) GetBalance(args *GetBalanceArgs, reply *GetBalanceReply) error {
	if !wallet.ValidateAddress(*args) {
		return errors.New("Address is not valid")
	}
	set := core.GetUTXOSet()

	*reply = 0
	pubKeyHash := utils.Base58Decode([]byte(*args))
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-4]
	utxos := set.FindUTXOByHash(pubKeyHash)

	for _, utxo := range utxos {
		*reply += utxo.Value
	}

	return nil
}

type GetBlocksArgs = struct {
	From int32
	To   int32
}
type GetBlocksReply = [][]byte

func GetBlocks(start, end int32) []*core.Block {
	var blocks GetBlocksReply
	log.Err(network.Call("DB.GetBlocks", &GetBlocksArgs{start, end}, &blocks))

	var ret []*core.Block
	for _, blockBytes := range blocks {
		ret = append(ret, core.BytesToBlock(blockBytes))
	}
	return ret
}

func (db *DB) GetBlocks(args *GetBlocksArgs, reply *GetBlocksReply) error {
	set := core.GetUTXOSet()
	log.Infof("GetBlocks args=%+v\n", args)
	for i := args.From; i <= args.To; i++ {
		data := set.Get(i)
		if data == nil {
			break
		}
		*reply = append(*reply, data)
	}
	return nil
}

func SendTransaction(txn *core.Transaction) {
	log.NotImplement()
}
