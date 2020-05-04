package api

import (
	"fmt"

	"github.com/YouDad/blockchain/core"
	"github.com/YouDad/blockchain/global"
	"github.com/YouDad/blockchain/network"
)

type ServerController struct {
	BaseController
}

type SendCMDArgs struct {
	SendFrom string
	SendTo   string
	Amount   int64
}

func SendCMD(from, to string, amount int64) error {
	args := SendCMDArgs{from, to, amount}
	return network.CallSelf("server/SendCMD", &args, nil)
}

func (c *ServerController) SendCMD() {
	var args SendCMDArgs
	c.ParseParameter(&args)

	set := core.GetUTXOSet(global.GetGroupByAddress(args.SendFrom))
	txn, err := set.CreateTransaction(args.SendFrom, args.SendTo, args.Amount)
	c.ReturnErr(err)
	c.ReturnErr(network.GetKnownNodes())
	GossipTxn(global.GetGroupByAddress(args.SendFrom), *txn,
		fmt.Sprintf("127.0.0.1:%s", global.Port))
	c.Return(nil)
}
