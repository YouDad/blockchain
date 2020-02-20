package main

import (
	"fmt"
	"testing"

	"github.com/YouDad/blockchain/app"
	"github.com/YouDad/blockchain/app/coin/coin_core"
	"github.com/YouDad/blockchain/core"
	"github.com/YouDad/blockchain/rpc"
)

func TestKnownNode(t *testing.T) {
	rpc.MyGetKnownNodes()
}

func TestPrintChain(t *testing.T) {
	core.InitCore(core.Config{
		DatabaseFile: fmt.Sprintf("blockchain_%s.db", "9999"),
		WalletFile:   fmt.Sprintf("wallet_%s.dat", "9999"),
	})
	core.InitCore(core.Config{
		GetAppdata: func() app.App {
			return coin_core.GetCoinApp(nil)
		},
	})
	set := coin_core.NewUTXOSet()
	set.Blocks().Foreach(func(k, v []byte) (isContinue bool) {
		if v == nil {
			fmt.Printf("[%x]: nil\n", k)
		} else {
			fmt.Printf("[%x]: %+v\n", k, core.DeserializeBlock(v))
		}
		return true
	})
}
