package api

import (
	"errors"

	"github.com/YouDad/blockchain/core"
	"github.com/YouDad/blockchain/log"
	"github.com/YouDad/blockchain/utils"
	"github.com/YouDad/blockchain/wallet"
)

type DB struct {
	set *core.UTXOSet
}

type GetGenesisReply = *core.Block

func (api *DB) GetGenesis(args *NIL, reply *GetGenesisReply) error {
	log.NotImplement()
	log.Infoln("GetGenesis")
	genesis := api.set.GetGenesis()
	*reply = genesis
	return nil
}

type GetBalanceArgs = string
type GetBalanceReply = int64

func (api *DB) GetBalance(args *GetBalanceArgs, reply *GetBalanceReply) error {
	if !wallet.ValidateAddress(*args) {
		return errors.New("Address is not valid")
	}

	*reply = 0
	pubKeyHash := utils.Base58Decode([]byte(*args))
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-4]
	utxos := api.set.FindUTXOByHash(pubKeyHash)

	for _, utxo := range utxos {
		*reply += utxo.Value
	}

	return nil
}
