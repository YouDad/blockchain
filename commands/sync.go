package commands

import (
	"github.com/YouDad/blockchain/api"
	"github.com/YouDad/blockchain/core"
	"github.com/YouDad/blockchain/log"
	"github.com/YouDad/blockchain/network"
	"github.com/YouDad/blockchain/store"
	"github.com/spf13/cobra"
)

var SyncCmd = &cobra.Command{
	Use:   "sync",
	Short: "sync information from other node",
	Run: func(cmd *cobra.Command, args []string) {
		log.Infoln("Syncing", Port)
		err := network.GetKnownNodes()
		network.Register(Port)
		if err != nil {
			log.Warnln(err)
		}

		var bc *core.Blockchain
		var utxoSet *core.UTXOSet
		if !store.IsDatabaseExists() {
			genesis, err := api.GetGenesis()
			if err != nil {
				log.Errln(err)
			}
			bc = core.CreateBlockchainFromGenesis(genesis)
			utxoSet = core.GetUTXOSet()
			utxoSet.Reindex()
		} else {
			bc = core.GetBlockchain()
			utxoSet = core.GetUTXOSet()
		}

		genesis := bc.GetGenesis()
		nowHeight := bc.GetHeight()
		height, err := api.SendVersion(nowHeight, genesis.Hash)
		if err == api.RootHashDifferentError {
			log.Warnln(err)
		} else if err == api.VersionDifferentError {
			log.Warnln(err)
		} else if err != nil {
			log.Warnln(err)
		}

		if height > nowHeight {
			blocks := api.GetBlocks(nowHeight+1, height)
			for _, block := range blocks {
				bc.AddBlock(block)
				utxoSet.Update(block)
			}
		}
	},
}
