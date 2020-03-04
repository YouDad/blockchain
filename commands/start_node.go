package commands

import (
	"github.com/spf13/cobra"

	"github.com/YouDad/blockchain/log"
	"github.com/YouDad/blockchain/network"
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
		network.Register(Port)
		go func() {
			//TODO: mining
		}()
		go func() {
			<-network.ServerReady
			SyncCmd.Run(cmd, args)
		}()
		network.StartServer()
	},
}
