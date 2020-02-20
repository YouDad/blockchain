package rpc

import (
	"errors"
	"log"

	"github.com/YouDad/blockchain/app/coin/coin_core"
	"github.com/YouDad/blockchain/app/coin/wallet"
	"github.com/YouDad/blockchain/core"
	"github.com/YouDad/blockchain/utils"
)

type DB struct {
	u *coin_core.UTXOSet
}

/*
 * DB's remote function
 * Local function {GetGenesis} can call this function
 */
type GetGenesisReply = []byte

func (db *DB) GetGenesis(args *NIL, reply *GetGenesisReply) error {
	genesis := db.u.GetGenesis()
	log.Println("GetGenesis", reply, core.DeserializeBlock(genesis))
	*reply = genesis
	return nil
}

/*
 * DB's remote function
 * Local function {GetBlocks} can call this function
 */
type GetBlocksArgs struct {
	From int
	To   int
}
type GetBlocksReply = [][]byte

func (db *DB) GetBlocks(args *GetBlocksArgs, reply *GetBlocksReply) error {
	log.Printf("GetBlocks args=%+v\n", args)
	for i := args.From; i <= args.To; i++ {
		data := db.u.GetByInt(i)
		if data == nil {
			break
		}
		*reply = append(*reply, data)
	}
	return nil
}

/*
 * DB's remote function
 * Local function {GetBalance} can call this function
 */
type GetBalanceArgs = string
type GetBalanceReply = int

func (db *DB) GetBalance(args *GetBalanceArgs, reply *GetBalanceReply) error {
	if wallet.ValidateAddress(*args) {
		return errors.New("Address is not valid")
	}

	*reply = 0
	pubKeyHash := utils.Base58Decode([]byte(*args))

	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-4]
	UTXOs := db.u.FindUTXO(pubKeyHash)

	for _, out := range UTXOs {
		*reply += out.Value
	}

	return nil
}

/*
 * DB's remote function
 * Local function {GetTransactions} can call this function
 */
type GetTransactionsArgs = NIL
type GetTransactionsReply = []*coin_core.Transaction

func (db *DB) GetTransactions(args *GetTransactionsArgs, reply *GetTransactionsReply) error {
	*reply = getTransactions()
	return nil
}
