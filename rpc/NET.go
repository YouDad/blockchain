package rpc

import (
	"bytes"
	"log"

	coin_core "github.com/YouDad/blockchain/app/coin/core"
	"github.com/YouDad/blockchain/core"
)

type NET struct {
	u *coin_core.UTXOSet
}

/*
 * NET's gossip function
 * Local function {SendTransaction} can call this function
 */
type SendTransactionArgs = coin_core.Transaction
type SendTransactionReply = byte

func (net *NET) SendTransaction(args *SendTransactionArgs, reply *SendTransactionReply) error {
	log.Println("SendTransaction")
	if !isTransactionExists(args) {
		addTransactionToMempool(args)
		SendTransaction(args)
	}
	return nil
}

/*
 * NET's remote function
 * Local function {GetKnownNodes} can call this function
 */
type GetKnownNodesArgs = string
type GetKnownNodesReply = []string

func (net *NET) GetKnownNodes(args *GetBalanceArgs, reply *GetKnownNodesReply) error {
	log.Println("GetKnownNodes", *args)
	if *args != "" {
		addKnownNode(*args)
	}
	var nodes []string
	for node, _ := range knownNodes {
		nodes = append(nodes, node)
	}
	*reply = nodes
	log.Println("GetKnownNodes return", *reply)
	return nil
}

/*
 * NET's remote function
 * Local function {SendVersion} can call this function
 */
type Version struct {
	Version  int
	Height   int
	RootHash []byte
}
type SendVersionArgs = Version
type SendVersionReply = Version

func (net *NET) SendVersion(args *SendVersionArgs, reply *SendVersionReply) error {
	log.Println("SendVersion")
	genesis := core.DeserializeBlock(net.u.GetGenesis())
	height := net.u.GetBestHeight()
	*reply = Version{
		Version:  version,
		Height:   height,
		RootHash: genesis.Hash,
	}

	if bytes.Compare(args.RootHash, genesis.Hash) == 0 {
		if args.Version != version {
			log.Printf("GetVersion: %d, NowVersion: %d\n", args.Version, version)
		}

		if args.Height > height {
			log.Printf("GetHeight: %d, NowHeight: %d\n", args.Height, height)
			blocks := GetBlocks(height+1, args.Height)
			for _, block := range blocks {
				net.u.AddBlock(block)
				net.u.Update(block)
			}
		}
	}
	return nil
}

/*
 * NET's gossip function
 * Local function {SendBlock} can call this function
 */
type SendBlockArgs = core.Block
type SendBlockReply = NIL

func (net *NET) SendBlock(args *SendBlockArgs, reply *SendBlockReply) error {
	log.Println("SendBlock")

	lastestBlock := core.DeserializeBlock(net.u.GetLastest())
	genesisBlock := core.DeserializeBlock(net.u.GetGenesis())
	lastestHash := lastestBlock.Hash
	height := lastestBlock.Height

	if !net.u.Blocks().IsExist(args.Hash) {
		SendBlock(args)
		net.u.Blocks().Set(args.Hash, []byte{})
	}

	if args.Height > height+1 {
		bestHeight, err := SendVersion(height, genesisBlock.Hash)
		if err == RootHashDifferentError {
			// TODO
			log.Println("Failed:", err)
		} else if err == VersionDifferentError {
			// TODO
			log.Println("Failed:", err)
		} else if err != nil {
			log.Println("Failed:", err)
		} else {
			if bestHeight > height {
				blocks := GetBlocks(bestHeight+1, height)
				for _, block := range blocks {
					net.u.AddBlock(block)
					net.u.Update(block)
				}
			}
		}
	}

	if args.Height == net.u.GetBestHeight()+1 && bytes.Equal(args.PrevBlockHash, lastestHash) {
		net.u.SetLastest(args.Hash, args.Serialize())
	}
	return nil
}

type HeartBeatArgs = NIL
type HeartBeatReply = NIL

func (net *NET) HeartBeat(args *HeartBeatArgs, reply *HeartBeatReply) error {
	return nil
}

type MyGetKnownNodesArgs = NIL
type MyGetKnwonNodesReply = PositionSlice

func (net *NET) MyGetKnownNodes(args *MyGetKnownNodesArgs, reply *MyGetKnwonNodesReply) error {
	*reply = sortedNodes
	return nil
}
