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

func SendVersion(group int, nowHeight int32, rootHash, nowHash types.HashValue) (int32, error, string) {
	var reply SendVersionReply
	args := SendVersionArgs{
		Group:    group,
		Version:  Version,
		Height:   nowHeight,
		RootHash: rootHash,
		NowHash:  nowHash,
	}
	err, address := network.CallInnerGroup("version/SendVersion", &args, &reply)

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
	reply := SendVersionReply{Group: -1}

	err := network.CallSelf("version/SendVersion", &reply, &reply)
	return reply, err
}

// @router /SendVersion [post]
func (c *VersionController) SendVersion() {
	var args SendVersionArgs
	var reply SendVersionReply
	c.ParseParameter(&args)
	log.Debugf("SendVersion %+v\n", args)

	if args.Group == -1 {
		args.Group = global.GetGroup()
	}

	bc := core.GetBlockchain()
	genesis := bc.GetGenesis(args.Group)
	lastest := bc.GetLastest(args.Group)
	lastestHeight := lastest.Height
	reply = types.Version{
		Group:    args.Group,
		Version:  Version,
		Height:   lastestHeight,
		RootHash: genesis.Hash(),
		NowHash:  lastest.Hash(),
	}

	c.Return(reply)
}
