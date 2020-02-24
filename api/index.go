package api

import (
	"net/rpc"

	"github.com/YouDad/blockchain/core"
	"github.com/YouDad/blockchain/log"
)

type NIL = bool

var NULL = true

func Register() {
	err := rpc.Register(&DB{core.GetUTXOSet()})
	log.Err(err)
	err = rpc.Register(&NET{core.GetUTXOSet()})
	log.Err(err)
}
