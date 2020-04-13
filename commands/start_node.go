package commands

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/YouDad/blockchain/api"
	"github.com/YouDad/blockchain/core"
	"github.com/YouDad/blockchain/log"
	"github.com/YouDad/blockchain/p2p"
	"github.com/YouDad/blockchain/storage"
)

var (
	startNodeAddress string
)

func init() {
	StartNodeCmd.Flags().StringVar(&startNodeAddress, "address", "", "node's coin address")
}

var StartNodeCmd = &cobra.Command{
	Use:   "start_node",
	Short: "Start a node with ID specified in port.",
	Run: func(cmd *cobra.Command, args []string) {
		log.Infoln("Starting node", Port)
		p2p.Register(Port)
		storage.RegisterDatabase(fmt.Sprintf("blockchain%s.db", Port))
		go func() {
			<-p2p.ServerReady
			err := p2p.GetKnownNodes()
			if err != nil {
				log.Warnln(err)
			}

			var bc *core.Blockchain
			if !bc.IsExists() {
				genesis, err := api.GetGenesis()
				if err != nil {
					log.Errln(err)
				}
				bc = core.CreateBlockchainFromGenesis(genesis)
			} else {
				bc = core.GetBlockchain()
			}
			utxoSet := core.NewUTXOSet()

			genesis := bc.GetGenesis()
			nowHeight := bc.GetHeight()
			height, err := api.SendVersion(nowHeight, genesis.Hash())
			if err == api.RootHashDifferentError {
				// TODO
				log.Warnln(err)
			} else if err == api.VersionDifferentError {
				// TODO
				log.Warnln(err)
			} else if err != nil {
				log.Warnln(err)
			}

			if height > nowHeight {
				blocks := api.GetBlocks(nowHeight+1, height)
				for _, block := range blocks {
					bc.AddBlock(block)
				}
				utxoSet.Reindex()
			}
		}()
		p2p.StartServer(startNodeAddress)
	},
}
