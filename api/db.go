package api

import (
	"bytes"
	"errors"

	"github.com/YouDad/blockchain/core"
	"github.com/YouDad/blockchain/global"
	"github.com/YouDad/blockchain/log"
	"github.com/YouDad/blockchain/network"
	"github.com/YouDad/blockchain/types"
	"github.com/YouDad/blockchain/utils"
	"github.com/YouDad/blockchain/wallet"
)

type DBController struct {
	BaseController
}

type GetGenesisReply = types.Block

var (
	ErrNull = errors.New("null")
)

func GetGenesis() (*types.Block, error) {
	var reply GetGenesisReply

	err, _ := network.CallInnerGroup("db/GetGenesis", nil, &reply)
	if err != nil {
		return nil, err
	}
	return &reply, err
}

// @router /GetGenesis [post]
func (c *DBController) GetGenesis() {
	var reply GetGenesisReply

	if global.IsDatabaseExists() {
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

	err := network.CallSelf("db/GetBalance", &args, &reply)
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
	Blocks []*types.Block
}

var ErrNoBlock = errors.New("No Needed Hash Block")

func CallbackGetBlocks(start, end int32, hash types.HashValue, address string) ([]*types.Block, error) {
	args := GetBlocksArgs{start, end, hash}
	var reply GetBlocksReply

	err := network.CallBack(address, "db/GetBlocks", &args, &reply)

	return reply.Blocks, err
}

func GetBlocks(start, end int32, hash types.HashValue) []*types.Block {
	args := GetBlocksArgs{start, end, hash}
	var reply GetBlocksReply
	err, _ := network.CallInnerGroup("db/GetBlocks", &args, &reply)
	log.Warn(err)

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

	bc := core.GetBlockchain()
	block := core.BytesToBlock(bc.Get(args.From))
	if block == nil {
		c.ReturnErr(ErrNoBlock)
	}

	if bytes.Compare(block.PrevHash, args.Hash) != 0 {
		log.Warnf("%x != %x\n", block.PrevHash, args.Hash)
		log.Warnln(block)
		block := core.BytesToBlock(bc.Get(args.From - 1))
		log.Warnln(block)
		c.ReturnErr(ErrNoBlock)
	}
	for i := args.From; i <= args.To; i++ {
		data := bc.Get(i)
		if data == nil {
			break
		}
		reply.Blocks = append(reply.Blocks, core.BytesToBlock(data))
	}
	c.Return(reply)
}

type GossipTxnArgs = types.Transaction

func GossipTxn(txn *types.Transaction) error {
	return network.GossipCallInnerGroup("db/GossipTxn", txn, nil)
}

// @router /GossipTxn [post]
func (c *DBController) GossipTxn() {
	var args GossipTxnArgs
	c.ParseParameter(&args)

	mempool := global.GetMempool()
	if !mempool.IsTxnExists(args) {
		bc := core.GetBlockchain()
		if bc.VerifyTransaction(args) {
			mempool.AddTxnToMempool(args)
			go GossipTxn(&args)
		} else {
			log.Warnf("AddTxnToMempool Verify false %x\n", args.Hash)
		}
	}
	c.Return(nil)
}

type GossipBlockArgs = types.Block

func CallbackGossipBlock(block *types.Block, address string) {
	network.CallBack(address, "db/GossipBlock", block, nil)
}

func GossipBlock(block *types.Block) {
	network.GossipCallInnerGroup("db/GossipBlock", block, nil)
}

// @router /GossipBlock [post]
func (c *DBController) GossipBlock() {
	var args GossipBlockArgs
	c.ParseParameter(&args)

	bc := core.GetBlockchain()
	set := core.GetUTXOSet()
	lastest := bc.GetLastest()
	lastestHeight := lastest.Height

	log.Debugf("GossipBlock get=%d, lastest=%d\n", args.Height, lastestHeight)

	if args.Height < lastestHeight {
		CallbackGossipBlock(lastest, c.GetString("address"))
	}

	if args.Height == lastestHeight+1 {
		if bytes.Compare(args.PrevHash, lastest.Hash()) == 0 {
			bc.AddBlock(&args)
			set.Update(&args)
		}
	}
	SyncBlocks(args.Height, c.GetString("address"))

	c.Return(nil)
}

type GetHashArgs struct{ Height int32 }
type GetHashReply struct{ Hash types.HashValue }

func CallbackGetHash(height int32, address string) (types.HashValue, error) {
	args := GetHashArgs{height}
	var reply GetHashReply

	err := network.CallBack(address, "db/GetHash", &args, &reply)

	return reply.Hash, err
}

func GetHash(height int32) (types.HashValue, error) {
	args := GetHashArgs{height}
	var reply GetHashReply

	err, _ := network.CallInnerGroup("db/GetHash", &args, &reply)

	return reply.Hash, err
}

// @router /GetHash [post]
func (c *DBController) GetHash() {
	log.Infoln(log.Funcname(0), c.GetString("address"))
	var args GetHashArgs
	var reply GetHashReply
	c.ParseParameter(&args)

	bc := core.GetBlockchain()
	block := core.BytesToBlock(bc.Get(args.Height))
	if block == nil {
		c.ReturnErr(ErrNoBlock)
	}
	reply.Hash = block.Hash()

	c.Return(reply)
}
