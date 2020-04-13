package apicb

import (
	"github.com/YouDad/blockchain/core"
	"github.com/YouDad/blockchain/log"
)

type DBApi struct {
	*core.UTXOSet
}

type NIL = bool
type GetGenesisReply = core.Block

func (api *DBApi) GetGenesis(args *NIL, reply *GetGenesisReply) error {
	log.Errln("NotImplement")
	log.Infoln("GetGenesis")
	genesis := api.UTXOSet.GetGenesis()
	*reply = genesis
	return nil
}
