package network

import (
	"errors"
	"math/rand"
	"net/rpc"
	"runtime/debug"
	"time"

	"github.com/YouDad/blockchain/log"
	"github.com/YouDad/blockchain/utils"
)

func call(node, method string, args interface{}, reply interface{}) error {
	log.Infoln("Request", method, node)
	client, err := rpc.DialHTTP(protocol, node)
	if err != nil {
		return err
	}
	defer client.Close()

	if utils.InterfaceIsNil(args) || utils.InterfaceIsNil(reply) {
		log.Warnln("args", args, "reply", reply)
		debug.PrintStack()
	}
	return client.Call(method, args, reply)
}

func Call(method string, args interface{}, reply interface{}) error {
	log.Infoln("Call", method)
	for _, node := range GetSortedNodes() {
		err := call(node.Address, method, args, reply)
		if err != nil {
			log.Warnln(node.Address, err)
			continue
		}
		return nil
	}
	return errors.New("None of the nodes responded!")
}

func CallMySelf(method string, args interface{}, reply interface{}) error {
	log.Infoln("CallMySelf", method)
	return call("127.0.0.1:"+Port, method, args, reply)
}

func random(min, max int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Int()%(max-min) + min
}

func GossipCall(method string, args interface{}, reply interface{}) {
	log.Infof("GossipCall %s\n", method)
	len := len(sortedNodes)
	half := len / 2
	visit := make([]bool, len)

	canVisit := func(min, max int) bool {
		for _, v := range visit[min:max] {
			if !v {
				return true
			}
		}
		return false
	}

	send := func(min, max int) bool {
		for {
			if !canVisit(min, max) {
				return false
			}

			visitor := random(min, max)
			if visit[visitor] {
				continue
			}
			visit[visitor] = true
			err := call(sortedNodes[visitor].Address, method, args, reply)
			if err != nil {
				log.Infoln(sortedNodes[visitor].Address, err)
				continue
			}
			log.Infoln(sortedNodes[visitor].Address, "success!")
			return true
		}
	}

	visited := 0

	if !send(0, half) {
		if send(half, len) {
			visited++
		}
	} else {
		visited++
	}

	if !send(0, half) {
		if send(half, len) {
			visited++
		}
	} else {
		visited++
	}

	if send(half, len) {
		visited++
	}

	if visited == 0 {
		log.Warnln("None of the nodes responded!")
	}
}
