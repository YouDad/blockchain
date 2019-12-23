package commands

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"

	"github.com/YouDad/blockchain/app/coin/wallet"
	"github.com/YouDad/blockchain/p2p"
)

var (
	startNodeAddress string
)

func init() {
	StartNodeCmd.Flags().StringVar(&startNodeAddress, "address", "",
		"node's coin address")
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
		p2p.StartServer(Port, startNodeAddress)
	},
}
