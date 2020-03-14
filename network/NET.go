package network

import (
	"github.com/YouDad/blockchain/conf"
	"github.com/YouDad/blockchain/global"
	"github.com/YouDad/blockchain/log"
)

type NET struct{}

type HeartBeatArgs = struct {
	Address string
}

func heartBeat(address string) {
	args := HeartBeatArgs{address}

	call(address, "net/HeartBeat", &args, nil)
}

func (net *NET) HeartBeat(args *HeartBeatArgs, reply *conf.NIL) error {
	log.Debugln("HeartBeat from", args.Address)

	return nil
}

type GetKnownNodesArgs = struct {
	Address string
}
type GetKnownNodesReply = struct {
	Addresses []string
}

func getKnownNodes(myAddress string, knownNodeAddresses *[]string) error {
	args := GetKnownNodesArgs{myAddress}
	var reply GetKnownNodesReply

	err := Call("net/GetKnownNodes", &args, &reply)
	*knownNodeAddresses = reply.Addresses
	return err
}

func (net *NET) GetKnownNodes(args *GetKnownNodesArgs, reply *GetKnownNodesReply) error {
	log.Infoln("GetKnownNodes", *args)

	knownNodes := global.GetKnownNodes()
	if args.Address != "" {
		knownNodes.AddNode(args.Address)
	}

	var nodes []string
	defer knownNodes.Release()
	for node, _ := range knownNodes.Get() {
		nodes = append(nodes, node)
	}
	reply.Addresses = nodes

	return nil
}
