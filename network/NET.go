package network

import (
	"github.com/YouDad/blockchain/conf"
	"github.com/YouDad/blockchain/log"
)

type NET struct{}

type HeartBeatArgs string

func heartBeat(address string) {
	call(address, "NET.HeartBeat", &address, &conf.NULL)
}

func (net *NET) HeartBeat(args *HeartBeatArgs, reply *conf.NIL) error {
	log.Debugln("NET.HeartBeat from", *args)

	return nil
}

type GetKnownNodesArgs = string
type GetKnownNodesReply = []string

func getKnownNodes(myAddress GetKnownNodesArgs, knownNodeAddresses *GetKnownNodesReply) error {
	return Call("NET.GetKnownNodes", &myAddress, knownNodeAddresses)
}

func (net *NET) GetKnownNodes(args *GetKnownNodesArgs, reply *GetKnownNodesReply) error {
	log.Infoln("GetKnownNodes", *args)

	if *args != "" {
		addKnownNode(*args)
	}

	var nodes []string
	for node, _ := range knownNodes {
		nodes = append(nodes, node)
	}
	*reply = nodes

	return nil
}
