package api

import (
	"bytes"
	"errors"
	"fmt"

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

	if reply.Version != args.Version {
		err = errors.New(
			fmt.Sprintf("Version is different. Args: %d, Reply: %d",
				args.Version, reply.Version,
			),
		)
	}

	if bytes.Compare(reply.RootHash, args.RootHash) != 0 {
		err = errors.New(
			fmt.Sprintf("RootHash is different. Args: %s, Reply: %s",
				args.RootHash, reply.RootHash),
		)
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
	c.ParseParameter(&args)
	log.Debugf("SendVersionArgs: %+v\n", args)

	if args.Group == -1 {
		args.Group = global.GetGroup()
	}

	bc := core.GetBlockchain(args.Group)
	genesis := bc.GetGenesis()
	lastest := bc.GetLastest()
	if genesis == nil || lastest == nil {
		c.Return(types.Version{
			Group:   args.Group,
			Version: Version,
			Height:  -1,
		})
	}
	c.Return(types.Version{
		Group:    args.Group,
		Version:  Version,
		Height:   lastest.Height,
		RootHash: genesis.Hash(),
		NowHash:  lastest.Hash(),
	})
}
