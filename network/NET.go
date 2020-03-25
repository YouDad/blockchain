package network

import (
	"time"

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

func (net *NET) HeartBeat(args *HeartBeatArgs) error {
	log.Debugln("HeartBeat from", args.Address)

	return nil
}

type GetKnownNodesArgs = struct {
	Address   string
	Timestamp int64
	Groups    []int
}
type GetKnownNodesReply = struct {
	Addresses []GetKnownNodesArgs
}

func getKnownNodes(myAddress string, knownNodeAddresses *[]GetKnownNodesArgs) error {
	args := GetKnownNodesArgs{myAddress, time.Now().UnixNano(), []int{0}}
	var reply GetKnownNodesReply

	err, _ := CallInnerGroup("net/GetKnownNodes", &args, &reply)
	*knownNodeAddresses = reply.Addresses
	return err
}

func (net *NET) GetKnownNodes(args *GetKnownNodesArgs, reply *GetKnownNodesReply) error {
	log.Infoln("GetKnownNodes", *args)

	knownNodes := global.GetKnownNodes()
	if args.Address != "" {
		knownNodes.AddNode(args.Address, args.Timestamp, args.Groups)
	}

	var nodes []GetKnownNodesArgs
	defer knownNodes.Release()
	for address, node := range knownNodes.Get() {
		nodes = append(nodes, GetKnownNodesArgs{
			address, node.Timestamp, node.Groups,
		})
	}
	reply.Addresses = nodes

	return nil
}
