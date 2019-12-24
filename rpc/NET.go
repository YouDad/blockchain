package rpc

import (
	"bytes"
	"fmt"

	coin_core "github.com/YouDad/blockchain/app/coin/core"
	"github.com/YouDad/blockchain/core"
	"github.com/YouDad/blockchain/log"
)

type NET struct {
	u *coin_core.UTXOSet
}

/*
 * NET's remote function
 * Local function {SendTx} can call this function
 * TODO
 */
type SendTxArgs coin_core.Transaction
type SendTxReply byte

func (net *NET) SendTx(args *SendTxArgs, reply *SendTxReply) error {
	log.Println("SendTx")
	fmt.Println("Get a sendtx")
	return nil
}

/*
 * NET's remote function
 * Local function {GetKnownNodes} can call this function
 */
type GetKnownNodesReply = []string

func (net *NET) GetKnownNodes(args *NIL, reply *GetKnownNodesReply) error {
	log.Println("GetKnownNodes")
	*reply = make(GetKnownNodesReply, 0)
	for node, _ := range knownNodes {
		*reply = append(*reply, node)
	}
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
	genesis := core.DeserializeBlock(net.u.Blocks().GetGenesis())
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
			//FIXME
			log.Printf("GetHeight: %d, NowHeight: %d\n", args.Height, height)
		}
	}
	return nil
}
