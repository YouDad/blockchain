package rpc

import (
	"errors"
	"net/rpc"
	"runtime/debug"

	"github.com/YouDad/blockchain/log"
	"github.com/YouDad/blockchain/p2p"
	"github.com/YouDad/blockchain/utils"
)

const (
	protocol = "tcp"
)

func call(node, method string, args interface{}, reply interface{}) error {
	client, err := rpc.DialHTTP(protocol, node)
	if err != nil {
		return err
	}

	if utils.InterfaceIsNil(args) || utils.InterfaceIsNil(reply) {
		log.Warnln("args", args, "reply", reply)
		debug.PrintStack()
	}
	return client.Call(method, args, reply)
}

func Call(method string, args interface{}, reply interface{}) error {
	log.Infoln("Call", method)
	// for node, _ := range knownNodes {
	for _, node := range p2p.GetSortedNodes() {
		err := call(node.Address, method, args, reply)
		if err != nil {
			log.Warnln(node.Address, err)
			continue
		}
		return nil
	}
	return errors.New("None of the nodes responded!")
}
