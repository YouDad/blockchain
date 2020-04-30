package api

import (
	"bytes"
	"errors"
	"fmt"

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

type GetGenesisArgs = struct {
	Group int
}
type GetGenesisReply = types.Block

var (
	ErrNull = errors.New("null")
)

func GetGenesis(group int) (*types.Block, error) {
	args := GetGenesisArgs{group}
	var reply GetGenesisReply

	err, _ := network.CallInnerGroup("db/GetGenesis", &args, &reply)
	if err != nil {
		return nil, err
	}
	return &reply, err
}

// @router /GetGenesis [post]
func (c *DBController) GetGenesis() {
	var args GetGenesisArgs
	var reply GetGenesisReply
	c.ParseParameter(&args)

	bc := core.GetBlockchain(args.Group)
	genesis := bc.GetGenesis()
	if genesis == nil {
		c.ReturnErr(errors.New(fmt.Sprintf("Blockchain[%d] don't have genesis", args.Group)))
	}
	reply = *genesis
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
	set := core.GetUTXOSet(global.GetGroup())

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
	Group int
	From  int32
	To    int32
	Hash  types.HashValue
}
type GetBlocksReply = struct {
	Blocks []*types.Block
}

func CallbackGetBlocks(group int, start, end int32, hash types.HashValue, address string) ([]*types.Block, error) {
	args := GetBlocksArgs{group, start, end, hash}
	var reply GetBlocksReply

	err := network.CallBack(address, "db/GetBlocks", &args, &reply)

	return reply.Blocks, err
}

func GetBlocks(group int, start, end int32, hash types.HashValue) []*types.Block {
	args := GetBlocksArgs{group, start, end, hash}
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

	bc := core.GetBlockchain(args.Group)
	block := bc.GetBlockByHeight(args.From)
	if block == nil {
		c.ReturnErr(errors.New(fmt.Sprintf(
			"No Needed Hash Block, Blockchain[%d].%d is nil", args.Group, args.From)))
	}

	if bytes.Compare(block.PrevHash, args.Hash) != 0 {
		log.Warnf("%s != %s\n", block.PrevHash, args.Hash)
		log.Warnln(block)
		block := bc.GetBlockByHeight(args.From - 1)
		log.Warnln(block)
		c.ReturnErr(errors.New(fmt.Sprintf("No Needed Hash Block, Hash is different.")))
	}
	for i := args.From; i <= args.To; i++ {
		data := bc.GetBlockByHeight(i)
		if data == nil {
			break
		}
		reply.Blocks = append(reply.Blocks, data)
	}
	c.Return(reply)
}

type GossipTxnArgs = struct {
	Txn   types.Transaction
	Group int
}

func GossipTxn(group int, txn types.Transaction) error {
	return network.GossipCallInnerGroup("db/GossipTxn", &GossipTxnArgs{txn, group}, nil)
}

// @router /GossipTxn [post]
func (c *DBController) GossipTxn() {
	var args GossipTxnArgs
	c.ParseParameter(&args)
	if !utils.InGroup(args.Group, global.GetGroup(), global.GroupNum, global.MaxGroupNum) {
		c.Return(nil)
	}

	if _, err := global.GetMempool(args.Group).GetTxn(args.Txn.Hash()); err != nil {
		bc := core.GetBlockchain(args.Group)
		if bc.VerifyTransaction(args.Txn) {
			global.GetMempool(args.Group).AddTxn(args.Txn)
			go GossipTxn(args.Group, args.Txn)
		} else {
			log.Warnf("AddTxn Verify false %s\n", args.Txn.Hash())
		}
	}
	c.Return(nil)
}

type GossipRelayTxnArgs = struct {
	FromGroup       int
	ToGroup         int
	Height          int32
	RelayMerklePath []types.MerklePath
	Txn             types.Transaction
}

func GossipRelayTxn(fromGroup int, toGroup int, height int32,
	relayMerklePath []types.MerklePath, txn *types.Transaction) {

	network.GossipCallSpecialGroup("db/GossipRelayTxn", &GossipRelayTxnArgs{
		fromGroup, toGroup, height, relayMerklePath, *txn}, nil, toGroup)
}

// @router /GossipRelayTxn
func (c *DBController) GossipRelayTxn() {
	var args GossipRelayTxnArgs
	c.ParseParameter(&args)

	if _, err := global.GetMempool(args.ToGroup).GetTxn(args.Txn.Hash()); err != nil {
		block := core.GetBlockhead(args.FromGroup).GetBlockheadByHeight(args.Height)
		if args.Txn.RelayVerify(block.MerkleRoot, args.RelayMerklePath) {
			global.GetMempool(args.ToGroup).AddTxn(args.Txn)
			go GossipRelayTxn(args.FromGroup, args.ToGroup, args.Height, args.RelayMerklePath, &args.Txn)
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

func CallSelfBlock(block *types.Block) {
	network.CallSelf("db/GossipBlock", block, nil)
	go GossipBlockHead(block)
}

// @router /GossipBlock [post]
func (c *DBController) GossipBlock() {
	var args GossipBlockArgs
	c.ParseParameter(&args)
	if !utils.InGroup(args.Group, global.GetGroup(), global.GroupNum, global.MaxGroupNum) {
		c.Return(nil)
	}

	log.Debugln("GossipBlock", "{{{{{{{{")
	bc := core.GetBlockchain(args.Group)
	set := core.GetUTXOSet(args.Group)
	lastest := bc.GetLastest()

	var lastestHeight int32 = -1
	if lastest != nil {
		lastestHeight = lastest.Height
	}

	log.Debugf("GossipBlock[%d] get=%d, lastest=%d\n",
		args.Group, args.Height, lastestHeight)

	// 认为对方的区块不够新，反向广播
	if args.Height < lastestHeight {
		CallbackGossipBlock(lastest, c.GetString("address"))
	}

	// 从高度上来说可能是后继区块
	if args.Height == lastestHeight+1 {
		// 满足哈希链
		if bytes.Compare(args.PrevHash, lastest.Hash()) == 0 {
			global.SyncMutex.Lock()
			bc.AddBlock(&args)
			set.Update(&args)
			global.SyncMutex.Unlock()
			go GossipBlock(&args)
			lastestHeight += 1
		}
	}

	// 认为我们和主链差一些区块
	if args.Height > lastestHeight {
		SyncBlocks(args.Group, args.Height, c.GetString("address"))
	}

	log.Debugln("GossipBlock", "}}}}}}}}")
	c.Return(nil)
}

type GossipBlockHeadArgs = types.Block

func GossipBlockHead(block *types.Block) {
	txns := block.Txns
	block.Txns = nil
	network.GossipCallInterGroup("db/GossipBlockHead", block, nil)
	block.Txns = txns
}

// @router /GossipBlockHead [post]
func (c *DBController) GossipBlockHead() {
	var args GossipBlockHeadArgs
	c.ParseParameter(&args)
	if !utils.InGroup(args.Group, global.GetGroup(), global.GroupNum, global.MaxGroupNum) {
		c.Return(nil)
	}

	bh := core.GetBlockhead(args.Group)
	if args.Verify() {
		if bh.AddBlockhead(&args) {
			go GossipBlockHead(&args)
		}
	} else {
		log.Warnln("AddBlockhead Verify failed")
	}

	c.Return(nil)
}

type GetHashArgs struct {
	Group  int
	Height int32
}
type GetHashReply struct{ Hash types.HashValue }

func CallbackGetHash(group int, height int32, address string) (types.HashValue, error) {
	args := GetHashArgs{group, height}
	var reply GetHashReply

	err := network.CallBack(address, "db/GetHash", &args, &reply)

	return reply.Hash, err
}

func GetHash(group int, height int32) (types.HashValue, error) {
	args := GetHashArgs{group, height}
	var reply GetHashReply

	err, _ := network.CallInnerGroup("db/GetHash", &args, &reply)

	return reply.Hash, err
}

// @router /GetHash [post]
func (c *DBController) GetHash() {
	var args GetHashArgs
	var reply GetHashReply
	c.ParseParameter(&args)

	bc := core.GetBlockchain(args.Group)
	block := bc.GetBlockByHeight(args.Height)
	if block == nil {
		c.ReturnErr(errors.New(fmt.Sprintf(
			"No Needed Hash Block, Height is %d", args.Height)))
	}
	reply.Hash = block.Hash()

	c.Return(reply)
}
