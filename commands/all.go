package commands

import (
	"github.com/spf13/cobra"

	"github.com/YouDad/blockchain/log"
	"github.com/YouDad/blockchain/network"
)

var AllCmd = &cobra.Command{
	Use:   "all",
	Short: "Start an all node",
	Run: func(cmd *cobra.Command, args []string) {
		log.Infoln("Starting node", Port)
		network.Register(Port)
		go func() {
			<-network.ServerReady
			SyncCmd.Run(cmd, args)
		}()
		network.StartServer()
	},
}