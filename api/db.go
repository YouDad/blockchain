package api

import (
	"bytes"
	"errors"

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

type DBController struct {
	BaseController
}

type GetGenesisReply = core.Block

var (
	ErrNull = errors.New("null")
)

func GetGenesis() (*core.Block, error) {
	var reply GetGenesisReply

	err := network.Call("db/GetGenesis", nil, &reply)
	if err != nil {
		return nil, err
	}
	return &reply, err
}

// @router /GetGenesis [post]
func (c *DBController) GetGenesis() {
	var reply GetGenesisReply

	if store.IsDatabaseExists() {
		genesis := core.GetBlockchain().GetGenesis()
		reply = *genesis
	}

	c.Return(reply)
}

type GetBalanceArgs = struct {
	Address string
}
type GetBalanceReply = struct {
	Balance int64
}

func GetBalance(address string) (int64, error) {
	args := GetBalanceArgs{address}
	var reply GetBalanceReply

	err := network.CallMySelf("db/GetBalance", &args, &reply)
	return reply.Balance, err
}

// @router /GetBalance [post]
func (c *DBController) GetBalance() {
	var args GetBalanceArgs
	var reply GetBalanceReply
	c.ParseParameter(&args)

	if !wallet.ValidateAddress(args.Address) {
		c.ReturnJson(SimpleJSONResult{"Address is not valid", nil})
	}
	set := core.GetUTXOSet()

	reply.Balance = 0
	pubKeyHash := utils.Base58Decode([]byte(args.Address))
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-4]
	utxos := set.FindUTXOByHash(pubKeyHash)

	for _, utxo := range utxos {
		reply.Balance += utxo.Value
	}

	c.Return(reply)
}

type GetBlocksArgs = struct {
	From int32
	To   int32
	Hash types.HashValue
}
type GetBlocksReply = struct {
	Blocks []*core.Block
}

var ErrNoBlock = errors.New("No Needed Hash Block")

func CallbackGetBlocks(start, end int32, hash types.HashValue, address string) ([]*core.Block, error) {
	args := GetBlocksArgs{start, end, hash}
	var reply GetBlocksReply

	err := network.Callback(address, "db/GetBlocks", &args, &reply)

	return reply.Blocks, err
}

func GetBlocks(start, end int32, hash types.HashValue) []*core.Block {
	args := GetBlocksArgs{start, end, hash}
	var reply GetBlocksReply

	log.Warn(network.Call("db/GetBlocks", &args, &reply))

	return reply.Blocks
}

// @router /GetBlocks [post]
func (c *DBController) GetBlocks() {
	var args GetBlocksArgs
	var reply GetBlocksReply
	c.ParseParameter(&args)

	if args.From == 0 {
		c.ReturnErr(errors.New("Height 0 use GetGenesis"))
	}

	set := core.GetUTXOSet()
	block := core.BytesToBlock(set.SetTable(conf.BLOCKS).Get(args.From))
	if block == nil {
		c.ReturnErr(ErrNoBlock)
	}

	if bytes.Compare(block.PrevHash, args.Hash) != 0 {
		c.ReturnErr(ErrNoBlock)
	}
	for i := args.From; i <= args.To; i++ {
		data := set.SetTable(conf.BLOCKS).Get(i)
		if data == nil {
			break
		}
		reply.Blocks = append(reply.Blocks, core.BytesToBlock(data))
	}
	c.Return(reply)
}

type SendTransactionArgs = core.Transaction

func SendTransaction(txn *core.Transaction) {
	network.GossipCall("db/SendTransaction", txn, nil)
}

// @router /SendTransaction [post]
func (c *DBController) SendTransaction() {
	var args SendTransactionArgs
	c.ParseParameter(&args)

	if !mempool.IsTxnExists(args) {
		bc := core.GetBlockchain()
		if bc.VerifyTransaction(args) {
			mempool.AddTxnToMempool(args)
			go SendTransaction(&args)
		} else {
			log.Warnf("AddTxnToMempool Verify false %x\n", args.Hash)
		}
	}
	c.Return(nil)
}

type SendBlockArgs = core.Block

func SendBlock(block *core.Block) {
	network.GossipCall("db/SendBlock", block, nil)
}

// @router /SendBlock [post]
func (c *DBController) SendBlock() {
	var args SendBlockArgs
	c.ParseParameter(&args)

	set := core.GetUTXOSet()
	lastest := set.GetLastest()
	lastestHeight := lastest.Height

	log.Debugf("SendBlock get=%d, lastest=%d\n", args.Height, lastestHeight)

	if args.Height == lastestHeight+1 {
		if bytes.Compare(args.PrevHash, lastest.Hash) == 0 {
			set.AddBlock(&args)
			set.Update(&args)
		}
	}
	syncBlocks(args.Height, c.GetString("address"))

	c.Return(nil)
}

type GetHashArgs struct{ Height int32 }
type GetHashReply struct{ Hash types.HashValue }

func CallbackGetHash(height int32, address string) (types.HashValue, error) {
	args := GetHashArgs{height}
	var reply GetHashReply

	err := network.Callback(address, "db/GetHash", &args, &reply)

	return reply.Hash, err
}

func GetHash(height int32) (types.HashValue, error) {
	args := GetHashArgs{height}
	var reply GetHashReply

	err := network.Call("db/GetHash", &args, &reply)

	return reply.Hash, err
}

// @router /GetHash [post]
func (c *DBController) GetHash() {
	log.Infoln(log.Funcname(), c.GetString("address"))
	var args GetHashArgs
	var reply GetHashReply
	c.ParseParameter(&args)

	bc := core.GetBlockchain()
	block := core.BytesToBlock(bc.SetTable(conf.BLOCKS).Get(args.Height))
	if block == nil {
		c.ReturnErr(ErrNoBlock)
	}
	reply.Hash = block.Hash

	c.Return(reply)
}
