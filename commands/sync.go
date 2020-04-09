package commands

import (
	"github.com/YouDad/blockchain/api"
	"github.com/YouDad/blockchain/core"
	"github.com/YouDad/blockchain/global"
	"github.com/YouDad/blockchain/log"
	"github.com/YouDad/blockchain/network"
	"github.com/spf13/cobra"
)

func init() {
	SyncCmd.Flags().IntVar(&global.GroupNum, "group", 1, "process group of number")
	SyncCmd.Flags().StringVar(&global.Address, "address", "", "address of node")
	SyncCmd.MarkFlagRequired("address")
}

var SyncCmd = &cobra.Command{
	Use:   "sync",
	Short: "sync information from other node",
	Run: func(cmd *cobra.Command, args []string) {
		log.Infoln("Syncing", global.Port)
		network.Register()
		log.Warn(network.GetKnownNodes())

		// XXX
		group := global.GetGroup()
		bc := core.GetBlockchain(group)

		// for
		if bc.GetHeight() < 0 {
			genesis, err := api.GetGenesis(group)
			log.Err(err)
			bc.AddBlock(genesis)
			core.GetUTXOSet(group).Reindex()
		}

		genesis := bc.GetGenesis()
		lastest := bc.GetLastest()
		height, err, address := api.SendVersion(group, lastest.Height, genesis.Hash(), lastest.Hash())
		if err == api.RootHashDifferentError {
			log.Warnln(err)
			return
		} else if err == api.VersionDifferentError {
			log.Warnln(err)
			return
		} else if err != nil {
			log.Warnln(err)
		}

		api.SyncBlocks(group, height, address)
	},
}
