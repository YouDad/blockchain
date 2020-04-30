package network

import (
	"time"

	"github.com/YouDad/blockchain/global"
	"github.com/YouDad/blockchain/log"
)

type NET struct{}

func heartBeat(address string) {
	call(address, "net/HeartBeat", nil, nil)
}

func (net *NET) HeartBeat() error {
	return nil
}

type GetKnownNodesArgs = struct {
	Address     string
	Timestamp   int64
	GroupBase   int
	GroupNumber int
}
type GetKnownNodesReply = struct {
	Addresses []GetKnownNodesArgs
}

func getKnownNodes(myAddress string, knownNodeAddresses *[]GetKnownNodesArgs) error {
	args := GetKnownNodesArgs{myAddress, time.Now().UnixNano(), global.GetGroup(), global.GroupNum}
	var reply GetKnownNodesReply

	err, _ := CallInterGroup("net/GetKnownNodes", &args, &reply)
	*knownNodeAddresses = reply.Addresses
	return err
}

func (net *NET) GetKnownNodes(args *GetKnownNodesArgs, reply *GetKnownNodesReply) error {
	knownNodes := global.GetKnownNodes()
	if args.Address != "" {
		knownNodes.AddNode(args.Address, args.Timestamp, args.GroupBase, args.GroupNumber)
	}

	var nodes []GetKnownNodesArgs
	defer knownNodes.Release()
	for address, node := range knownNodes.Get() {
		nodes = append(nodes, GetKnownNodesArgs{
			address, node.Timestamp, node.GroupBase, node.GroupNumber,
		})
	}
	reply.Addresses = nodes
	log.Debugln("GetKnownNodes Reply", *reply)

	return nil
}
