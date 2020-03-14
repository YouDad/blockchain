package api

import (
	"bytes"
	"errors"

	"github.com/YouDad/blockchain/conf"
	"github.com/YouDad/blockchain/core"
	"github.com/YouDad/blockchain/log"
	"github.com/YouDad/blockchain/network"
	"github.com/YouDad/blockchain/types"
)

type VersionController struct {
	BaseController
}

type SendVersionArgs = types.Version
type SendVersionReply = types.Version

var (
	RootHashDifferentError = errors.New("RootHash is different.")
	VersionDifferentError  = errors.New("Version is different.")
)

func SendVersion(nowHeight int32, rootHash, nowHash types.HashValue) (int32, error) {
	var reply SendVersionReply
	args := SendVersionArgs{
		Version:  conf.Version,
		Height:   nowHeight,
		RootHash: rootHash,
		NowHash:  nowHash,
	}
	err := network.Call("version/SendVersion", &args, &reply)

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
	var reply SendVersionReply

	err := network.CallMySelf("version/SendVersion", &reply, &reply)
	return reply, err
}

// @router /SendVersion [post]
func (c *VersionController) SendVersion() {
	var args SendVersionArgs
	var reply SendVersionReply
	c.ParseParameter(&args)

	set := core.GetUTXOSet()
	genesis := set.GetGenesis()
	lastest := set.GetLastest()
	lastestHeight := lastest.Height
	reply = types.Version{
		Version:  conf.Version,
		Height:   lastestHeight,
		RootHash: genesis.Hash,
		NowHash:  lastest.Hash,
	}

	if args.Height == 0 && len(args.RootHash) == 0 {
		c.Return(reply)
	}

	if bytes.Compare(args.RootHash, genesis.Hash) != 0 {
		c.ReturnErr(RootHashDifferentError)
	}

	if args.Version != conf.Version {
		log.Infof("GetVersion: %d, WeVersion: %d\n", args.Version, conf.Version)
		log.Warnln("Version Update")
		c.ReturnErr(VersionDifferentError)
	}

	syncBlocks(lastestHeight, c.GetString("address"))

	c.Return(reply)
}
