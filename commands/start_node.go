package commands

import (
	"github.com/spf13/cobra"

	"github.com/YouDad/blockchain/core"
	"github.com/YouDad/blockchain/log"
	"github.com/YouDad/blockchain/rpc"
	"github.com/YouDad/blockchain/wallet"
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

		// 是否挖矿，检验钱包
		if len(startNodeAddress) > 0 {
			if wallet.ValidateAddress(startNodeAddress) {
				log.Infoln("Mining is on. Address to receive rewards:", startNodeAddress)
			} else {
				log.Errln("Wrong miner address!")
			}
		}

		rpc.Init(Port)

		var bc *core.Blockchain
		if !core.IsBlockchainExists() {
			genesis, err := rpc.GetGenesis()
			if err != nil {
				log.Errln(err)
			}
			bc = core.CreateBlockchainFromGenesis(genesis)
		} else {
			bc = core.NewBlockchain()
		}
		utxoSet := core.NewUTXOSet()

		go func() {
			<-rpc.ServerReady
			err := rpc.GetKnownNodes()
			if err != nil {
				log.Warnln(err)
			}

			genesisBlock := core.DeserializeBlock(bc.GetGenesis())
			bestHeight := bc.GetBestHeight()
			height, err := rpc.SendVersion(bestHeight, genesisBlock.Hash)
			if err == rpc.RootHashDifferentError {
				// TODO
				log.Warnln(err)
			} else if err == rpc.VersionDifferentError {
				// TODO
				log.Warnln(err)
			} else if err != nil {
				log.Warnln(err)
			}

			if height > bestHeight {
				blocks := rpc.GetBlocks(bestHeight+1, height)
				for _, block := range blocks {
					bc.AddBlock(block)
				}
				utxoSet.Reindex()
			}

			rpc.GetTransactions()

			bc.Close()
			utxoSet.Close()
		}()

		rpc.StartServer(Port, startNodeAddress)
	},
}
