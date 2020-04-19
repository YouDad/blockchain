package api

import (
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
	set := core.GetUTXOSet(global.GetGroup())
	txn, err := set.CreateTransaction(args.SendFrom, args.SendTo, args.Amount)
	c.ReturnErr(err)
	c.ReturnErr(network.GetKnownNodes())
	c.ReturnErr(GossipTxn(global.GetGroup(), *txn))
	c.Return(nil)
}
