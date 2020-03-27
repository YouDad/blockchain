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
	SyncCmd.Flags().UintVar(&global.GroupNum, "group", 1, "process group of number")
}

var SyncCmd = &cobra.Command{
	Use:   "sync",
	Short: "sync information from other node",
	Run: func(cmd *cobra.Command, args []string) {
		log.Infoln("Syncing", global.Port)
		network.Register()
		log.Warn(network.GetKnownNodes())

		// XXX
		bc := core.GetBlockchain()

		// for
		group := global.GetGroup()
		if bc.GetHeight(group) < 0 {
			genesis, err := api.GetGenesis(group)
			log.Err(err)
			bc.AddBlock(group, genesis)
			core.GetUTXOSet().Reindex(group)
		}

		genesis := bc.GetGenesis(group)
		lastest := bc.GetLastest(group)
		lastestHeight := lastest.Height
		lastestHash := lastest.Hash()
		height, err, address := api.SendVersion(group, lastestHeight, genesis.Hash(), lastestHash)
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
