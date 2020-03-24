package api

import (
	"bytes"
	"errors"

	"github.com/YouDad/blockchain/core"
	"github.com/YouDad/blockchain/global"
	"github.com/YouDad/blockchain/log"
	"github.com/YouDad/blockchain/network"
	"github.com/YouDad/blockchain/types"
)

type VersionController struct {
	BaseController
}

type SendVersionArgs = types.Version
type SendVersionReply = types.Version

const Version int = 0x00

var (
	RootHashDifferentError = errors.New("RootHash is different.")
	VersionDifferentError  = errors.New("Version is different.")
)

func SendVersion(nowHeight int32, rootHash, nowHash types.HashValue) (int32, error, string) {
	var reply SendVersionReply
	args := SendVersionArgs{
		Version:  Version,
		Height:   nowHeight,
		RootHash: rootHash,
		NowHash:  nowHash,
	}
	err, address := network.Call("version/SendVersion", &args, &reply)

	if err != nil {
		return 0, err, address
	}

	if reply.Version != Version {
		err = VersionDifferentError
	}

	if bytes.Compare(reply.RootHash, rootHash) != 0 {
		err = RootHashDifferentError
	}

	return reply.Height, err, address
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
	log.Debugf("SendVersion %+v\n", args)

	bc := core.GetBlockchain()
	genesis := bc.GetGenesis()
	lastest := bc.GetLastest()
	lastestHeight := global.GetHeight()
	reply = types.Version{
		Version:  Version,
		Height:   lastestHeight,
		RootHash: genesis.Hash(),
		NowHash:  lastest.Hash(),
	}

	if args.Height == 0 && len(args.RootHash) == 0 {
		c.Return(reply)
	}

	if bytes.Compare(args.RootHash, genesis.Hash()) != 0 {
		c.ReturnErr(RootHashDifferentError)
	}

	if args.Version != Version {
		log.Infof("GetVersion: %d, WeVersion: %d\n", args.Version, Version)
		log.Warnln("Version Update")
		c.ReturnErr(VersionDifferentError)
	}

	SyncBlocks(lastestHeight, c.GetString("address"))

	c.Return(reply)
}
