package api

import (
	"github.com/YouDad/blockchain/network"
	"github.com/YouDad/blockchain/types"
)

func GetVersion() (types.Version, error) {
	var version SendVersionReply
	err := network.CallMySelf("NET.SendVersion", &version, &version)
	return version, err
}
