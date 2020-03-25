package api

import (
	"github.com/YouDad/blockchain/core"
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
	set := core.GetUTXOSet()
	txn, err := set.NewUTXOTransaction(args.SendFrom, args.SendTo, args.Amount)
	if err != nil {
		c.ReturnErr(err)
	}
	err = network.GetKnownNodes()
	if err != nil {
		c.ReturnErr(err)
	}
	err = GossipTxn(txn)
	c.ReturnErr(err)
}
