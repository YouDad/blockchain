package api

import (
	"github.com/YouDad/blockchain/network"
)

type NetController struct {
	BaseController
	net network.NET
}

func (c *NetController) HeartBeat() {
	c.ParseParameter(nil)

	c.ReturnErr(c.net.HeartBeat())
	c.Return(nil)
}

func (c *NetController) GetKnownNodes() {
	var args network.GetKnownNodesArgs
	var reply network.GetKnownNodesReply
	c.ParseParameter(&args)

	c.ReturnErr(c.net.GetKnownNodes(&args, &reply))
	c.Return(reply)
}
