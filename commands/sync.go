package commands

import (
	"bytes"

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
		network.Register(Port)
		log.Warn(network.GetKnownNodes())

		var bc *core.Blockchain
		var utxoSet *core.UTXOSet
		if !store.IsDatabaseExists() {
			genesis, err := api.GetGenesis()
			log.Err(err)
			bc = core.CreateBlockchainFromGenesis(genesis)
			utxoSet = core.GetUTXOSet()
			utxoSet.Reindex()
		} else {
			bc = core.GetBlockchain()
			utxoSet = core.GetUTXOSet()
		}

		genesis := bc.GetGenesis()
		lastest := bc.GetLastest()
		lastestHeight := lastest.Height
		lastestHash := lastest.Hash
		height, err := api.SendVersion(lastestHeight, genesis.Hash)
		if err == api.RootHashDifferentError {
			log.Warnln(err)
		} else if err == api.VersionDifferentError {
			log.Warnln(err)
		} else if err != nil {
			log.Warnln(err)
		}

		if height > lastestHeight {
			blocks := api.GetBlocks(lastestHeight+1, height, lastest.Hash)
			for _, block := range blocks {
				if bytes.Compare(block.PrevHash, lastestHash) == 0 {
					bc.AddBlock(block)
					utxoSet.Update(block)
					lastestHash = block.Hash
				}
			}
		}
	},
}
