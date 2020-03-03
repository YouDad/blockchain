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
	height := set.GetHeight()
	*reply = types.Version{
		Version:  conf.Version,
		Height:   height,
		RootHash: genesis.Hash(),
	}

	if bytes.Compare(args.RootHash, genesis.Hash()) == 0 {
		if args.Version != conf.Version {
			log.Infof("GetVersion: %d, WeVersion: %d\n", args.Version, conf.Version)
			log.Warnln("Version Update") // TODO
		}

		if args.Height > height {
			log.Infof("GetHeight: %d, WeHeight: %d\n", args.Height, height)
			blocks := GetBlocks(height+1, args.Height)
			for _, block := range blocks {
				set.AddBlock(block)
				set.Update(block)
			}
		}
	}
	return nil
}
