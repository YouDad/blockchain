package api

import (
	"bytes"

	"github.com/YouDad/blockchain/core"
	"github.com/YouDad/blockchain/log"
	"github.com/YouDad/blockchain/types"
)

var version int = 0x00

type NET struct {
	set *core.UTXOSet
}

func (net *NET) HeartBeat(args *NIL, reply *NIL) error {
	return nil
}

type SendVersionArgs = types.Version
type SendVersionReply = types.Version

func (net *NET) SendVersion(args *SendVersionArgs, reply *SendVersionReply) error {
	log.Infoln("SendVersion")
	genesis := net.set.GetGenesis()
	height := net.set.GetHeight()
	*reply = types.Version{
		Version:  version,
		Height:   height,
		RootHash: genesis.Hash(),
	}

	if bytes.Compare(args.RootHash, genesis.Hash()) == 0 {
		if args.Version != version {
			log.Infof("GetVersion: %d, WeVersion: %d\n", args.Version, version)
			log.Warnln("Version Update") // TODO
		}

		if args.Height > height {
			log.Infof("GetHeight: %d, WeHeight: %d\n", args.Height, height)
			blocks := GetBlocks(height+1, args.Height)
			for _, block := range blocks {
				net.set.AddBlock(block)
				net.set.Update(block)
			}
		}
	}
	return nil
}
