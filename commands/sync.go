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

		var bc *core.Blockchain
		if !global.IsDatabaseExists() {
			genesis, err := api.GetGenesis()
			log.Err(err)
			bc = core.CreateBlockchainFromGenesis(genesis)
		} else {
			bc = core.GetBlockchain()
		}

		genesis := bc.GetGenesis()
		lastest := bc.GetLastest()
		lastestHeight := lastest.Height
		lastestHash := lastest.Hash()
		height, err, address := api.SendVersion(lastestHeight, genesis.Hash(), lastestHash)
		if err == api.RootHashDifferentError {
			log.Warnln(err)
			return
		} else if err == api.VersionDifferentError {
			log.Warnln(err)
			return
		} else if err != nil {
			log.Warnln(err)
		}

		api.SyncBlocks(height, address)
	},
}
