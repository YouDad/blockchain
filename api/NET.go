package api

import (
	"net/rpc"

	"github.com/YouDad/blockchain/log"
	"github.com/YouDad/blockchain/network"
)

type NET struct {
	network.NET
}

func init() {
	log.Err(rpc.Register(new(NET)))
}
