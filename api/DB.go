package api

import (
	"bytes"
	"errors"
	"net/rpc"

	"github.com/YouDad/blockchain/conf"
	"github.com/YouDad/blockchain/core"
	"github.com/YouDad/blockchain/log"
	"github.com/YouDad/blockchain/mempool"
	"github.com/YouDad/blockchain/network"
	"github.com/YouDad/blockchain/store"
	"github.com/YouDad/blockchain/types"
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
	Hash types.HashValue
}
type GetBlocksReply = [][]byte

var ErrNoBlock = errors.New("No Needed Hash Block")

func GetBlocks(start, end int32, hash types.HashValue) []*core.Block {
	var blocks GetBlocksReply
	log.Warn(network.Call("DB.GetBlocks", &GetBlocksArgs{start, end, hash}, &blocks))

	var ret []*core.Block
	for _, blockBytes := range blocks {
		ret = append(ret, core.BytesToBlock(blockBytes))
	}
	return ret
}

func (db *DB) GetBlocks(args *GetBlocksArgs, reply *GetBlocksReply) error {
	set := core.GetUTXOSet()
	log.Infof("GetBlocks args=%+v\n", args)
	block := core.BytesToBlock(set.Get(args.From))
	if bytes.Compare(block.PrevHash, args.Hash) != 0 {
		return ErrNoBlock
	}
	for i := args.From; i <= args.To; i++ {
		data := set.Get(i)
		if data == nil {
			break
		}
		*reply = append(*reply, data)
	}
	return nil
}

type SendTransactionArgs = core.Transaction
type SendTransactionReply = conf.NIL

func SendTransaction(txn *core.Transaction) {
	network.GossipCall("DB.SendTransaction", txn, &conf.NULL)
}

func (db *DB) SendTransaction(args *SendTransactionArgs, reply *SendTransactionReply) error {
	log.Infoln("SendTransaction")

	if !mempool.IsTxnExists(args) {
		bc := core.GetBlockchain()
		if bc.VerifyTransaction(args) {
			mempool.AddTxnToMempool(args)
			go network.GossipCall("DB.SendTransaction", args, &conf.NULL)
		} else {
			log.Warnf("AddTxnToMempool Verify false %x\n", args.Hash)
		}
	}
	return nil
}

type SendBlockArgs = core.Block
type SendBlockReply = conf.NIL

func SendBlock(block *SendBlockArgs) {
	network.GossipCall("DB.SendBlock", block, &conf.NULL)
}

func (db *DB) SendBlock(args *SendBlockArgs, reply *SendBlockReply) error {
	log.Infoln("SendBlock")

	utxoSet := core.GetUTXOSet()
	lastest := utxoSet.GetLastest()
	lastestHeight := lastest.Height
	lastestHash := lastest.Hash

	if args.Height > lastestHeight+1 {
		blocks := GetBlocks(lastestHeight+1, args.Height-1, lastestHash)
		for _, block := range blocks {
			if bytes.Compare(block.PrevHash, lastestHash) == 0 {
				utxoSet.AddBlock(block)
				utxoSet.Update(block)
				lastestHash = block.Hash
			} else {
				break
			}
		}
	}

	if args.Height == lastest.Height+1 {
		if bytes.Compare(args.PrevHash, lastestHash) == 0 {
			utxoSet.AddBlock(args)
			go network.GossipCall("DB.SendBlock", args, &conf.NULL)
		}
	}
	return nil
}
