package api

import (
	"github.com/YouDad/blockchain/network"
)

type NetController struct {
	BaseController
	net network.NET
}

func (c *NetController) HeartBeat() {
	var args network.HeartBeatArgs
	c.ParseParameter(&args)

	c.net.HeartBeat(&args, nil)
	c.Return(nil)
}

func (c *NetController) GetKnownNodes() {
	var args network.GetKnownNodesArgs
	var reply network.GetKnownNodesReply
	c.ParseParameter(&args)

	c.net.GetKnownNodes(&args, &reply)
	c.Return(reply)
}
