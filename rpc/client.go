package rpc

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"net/rpc"
	"runtime/debug"

	coin_core "github.com/YouDad/blockchain/app/coin/core"
	"github.com/YouDad/blockchain/core"
	"github.com/YouDad/blockchain/utils"
)

type NIL = bool

func call(node, method string, args interface{}, reply interface{}) error {
	client, err := rpc.DialHTTP(protocol, node)
	if err != nil {
		return err
	}

	if utils.InterfaceIsNil(args) || utils.InterfaceIsNil(reply) {
		log.Fatalln("args", args, "reply", reply)
		debug.PrintStack()
	}
	return client.Call(method, args, reply)
}

func Call(method string, args interface{}, reply interface{}) error {
	log.Println("Call", method)
	// for node, _ := range knownNodes {
	for _, node := range sortedNodes {
		err := call(node.Address, method, args, reply)
		if err != nil {
			log.Println("Failed:", node.Address, err)
			continue
		}
		return nil
	}
	return errors.New("None of the nodes responded!")
}

func CallMySelf(port string, method string, args interface{}, reply interface{}) error {
	log.Println("CallMySelf", method)
	return call(fmt.Sprintf("127.0.0.1:%s", port), method, args, reply)
}

func GetBalance(port string, address string) (int, error) {
	var balance int
	err := CallMySelf(port, "DB.GetBalance", &address, &balance)
	return balance, err
}

func GetVersion(port string) (Version, error) {
	var version SendVersionReply
	err := CallMySelf(port, "NET.SendVersion", &version, &version)
	return version, err
}

// Send New Transaction To Known Node
func SendTransaction(tx *coin_core.Transaction) {
	BOOL := true
	GossipCall("NET.SendTransaction", tx, &BOOL)
}

// Send New Block To Known Node
func SendBlock(block *core.Block) {
	BOOL := true
	log.Printf("{TMP} NET.SendBlock block=%+v\n", block)
	GossipCall("NET.SendBlock", block, &BOOL)
}

// Get Other Nodes From Known Node
func GetKnownNodes() error {
	knownNodeAddresses := new([]string)
	myAddress := ""
	myAddress = getExternIP() + ":" + Port
	err := Call("NET.GetKnownNodes", &myAddress, knownNodeAddresses)
	log.Printf("GetKnownNodes [%s] get:%+v", myAddress, *knownNodeAddresses)
	if err == nil {
		for _, knownNodeAddress := range *knownNodeAddresses {
			addKnownNode(knownNodeAddress)
		}
	}
	return err
}

var (
	RootHashDifferentError = errors.New("RootHash is different.")
	VersionDifferentError  = errors.New("Version is different.")
)

// Tell Known Node My Version Infomation To Get New Information
func SendVersion(height int, rootHash []byte) (int, error) {
	var reply SendVersionReply
	err := Call("NET.SendVersion", &SendVersionArgs{
		Version:  version,
		Height:   height,
		RootHash: rootHash,
	}, &reply)

	if err != nil {
		return 0, err
	}

	if reply.Version != version {
		err = VersionDifferentError
	}

	if bytes.Compare(reply.RootHash, rootHash) != 0 {
		err = RootHashDifferentError
	}

	return reply.Height, err
}

// From Known Node Get Genesis Block
func GetGenesis() (*core.Block, error) {
	var genesisBlock GetGenesisReply
	BOOL := true
	err := Call("DB.GetGenesis", &BOOL, &genesisBlock)
	if err != nil {
		return nil, err
	}
	return core.DeserializeBlock(genesisBlock), err
}

// From Known Node Get Lastest Blocks
func GetBlocks(from, to int) []*core.Block {
	var blocks GetBlocksReply
	err := Call("DB.GetBlocks", &GetBlocksArgs{from, to}, &blocks)
	if err != nil {
		log.Panic(err)
	}

	var ret []*core.Block
	for _, block := range blocks {
		ret = append(ret, core.DeserializeBlock(block))
	}
	return ret
}

// From Known Node Get Lastest Transactions
func GetTransactions() {
	var txs []*coin_core.Transaction
	BOOL := true
	err := Call("DB.GetTransactions", &BOOL, &txs)
	if err != nil {
		log.Println("Failed:", err)
	} else {
		for _, tx := range txs {
			addTransactionToMempool(tx)
		}
	}
}

func HeartBeat(address string) {
	BOOL := true
	call(address, "NET.HeartBeat", &BOOL, &BOOL)
}

func MyGetKnownNodes() {
	var nodes MyGetKnwonNodesReply
	BOOL := true
	err := CallMySelf("9999", "NET.MyGetKnownNodes", &BOOL, &nodes)
	if err != nil {
		log.Println("Failed:", err)
	} else {
		for _, node := range nodes {
			log.Printf("%+v\n", node)
		}
	}
}
