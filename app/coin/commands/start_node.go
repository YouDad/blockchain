package commands

import (
	"fmt"

	"github.com/spf13/cobra"

	coin_core "github.com/YouDad/blockchain/app/coin/core"
	"github.com/YouDad/blockchain/app/coin/wallet"
	"github.com/YouDad/blockchain/core"
	"github.com/YouDad/blockchain/log"
	"github.com/YouDad/blockchain/rpc"
)

var (
	startNodeAddress string
	startNodeIP      string
)

func init() {
	StartNodeCmd.Flags().StringVar(&startNodeAddress, "address", "",
		"node's coin address")
	StartNodeCmd.Flags().StringVar(&startNodeIP, "ip", "",
		"node's ip")
}

var StartNodeCmd = &cobra.Command{
	Use:   "start_node",
	Short: "Start a node with ID specified in port.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Starting node %s\n", Port)
		if len(startNodeAddress) > 0 {
			if wallet.ValidateAddress(startNodeAddress) {
				fmt.Println("Mining is on. Address to receive rewards: ", startNodeAddress)
			} else {
				log.Panic("Wrong miner address!")
			}
		}

		var bc *core.Blockchain
		if !core.IsBlockchainExists() {
			genesis, err := rpc.GetGenesis()
			if err != nil {
				log.Panic(err)
			}
			bc = core.CreateBlockchainFromGenesis(genesis)
		} else {
			bc = core.NewBlockchain()
		}
		utxo_set := coin_core.NewUTXOSet()

		err := rpc.GetKnownNodes()
		if err != nil {
			log.Println(err)
		}

		genesisBlock := core.DeserializeBlock(bc.GetGenesis())
		bestHeight := bc.GetBestHeight()
		height, err := rpc.SendVersion(bestHeight, genesisBlock.Hash)
		if err == rpc.RootHashDifferentError {
			// TODO
			log.Println(err)
		} else if err == rpc.VersionDifferentError {
			// TODO
			log.Println(err)
		} else if err != nil {
			log.Println(err)
		}

		if height > bestHeight {
			blocks := rpc.GetBlocks(bestHeight+1, height)
			for _, block := range blocks {
				bc.AddBlock(block)
			}
			utxo_set.Reindex()
		}

		err = rpc.GetTransactions()
		if err != nil {
			log.Println(err)
		}

		bc.Close()
		utxo_set.Close()

		rpc.StartServer(Port, startNodeAddress)
	},
}
