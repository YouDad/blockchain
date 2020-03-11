package api

import (
	"bytes"
	"errors"
	"net/rpc"

	"github.com/YouDad/blockchain/conf"
	"github.com/YouDad/blockchain/core"
	"github.com/YouDad/blockchain/log"
	"github.com/YouDad/blockchain/network"
	"github.com/YouDad/blockchain/types"
)

type VERSION struct{}

func init() {
	log.Err(rpc.Register(new(VERSION)))
}

type SendVersionArgs = types.Version
type SendVersionReply = types.Version

var (
	RootHashDifferentError = errors.New("RootHash is different.")
	VersionDifferentError  = errors.New("Version is different.")
)

func SendVersion(nowHeight int32, rootHash types.HashValue) (int32, error) {
	var reply SendVersionReply
	err := network.Call("VERSION.SendVersion", &SendVersionArgs{
		Version:  conf.Version,
		Height:   nowHeight,
		RootHash: rootHash,
	}, &reply)

	if err != nil {
		return 0, err
	}

	if reply.Version != conf.Version {
		err = VersionDifferentError
	}

	if bytes.Compare(reply.RootHash, rootHash) != 0 {
		err = RootHashDifferentError
	}

	return reply.Height, err
}

func GetVersion() (types.Version, error) {
	var version SendVersionReply
	err := network.CallMySelf("VERSION.SendVersion", &version, &version)
	return version, err
}

func (ver *VERSION) SendVersion(args *SendVersionArgs, reply *SendVersionReply) error {
	set := core.GetUTXOSet()
	log.Infoln("SendVersion")
	genesis := set.GetGenesis()
	lastest := set.GetLastest()
	lastestHeight := lastest.Height
	lastestHash := lastest.Hash
	*reply = types.Version{
		Version:  conf.Version,
		Height:   lastestHeight,
		RootHash: genesis.Hash,
	}

	if bytes.Compare(args.RootHash, genesis.Hash) == 0 {
		if args.Version != conf.Version {
			log.Infof("GetVersion: %d, WeVersion: %d\n", args.Version, conf.Version)
			log.Warnln("Version Update") // TODO
		}

		if args.Height > lastestHeight {
			log.Infof("GetHeight: %d, WeHeight: %d\n", args.Height, lastestHeight)
			blocks := GetBlocks(lastestHeight+1, args.Height, lastestHash)
			for _, block := range blocks {
				if bytes.Compare(block.PrevHash, lastestHash) == 0 {
					set.AddBlock(block)
					set.Update(block)
					lastestHash = block.Hash
				} else {
					break
				}
			}
		}
	}
	return nil
}
