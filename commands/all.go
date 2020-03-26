package commands

import (
	"github.com/spf13/cobra"

	"github.com/YouDad/blockchain/global"
	"github.com/YouDad/blockchain/log"
	"github.com/YouDad/blockchain/network"
)

func init() {
	AllCmd.Flags().UintVar(&global.GroupNum, "group", 1, "process group of number")
}

var AllCmd = &cobra.Command{
	Use:   "all",
	Short: "Start an all node",
	Run: func(cmd *cobra.Command, args []string) {
		log.Infoln("Starting node", global.Port)
		network.Register()
		go func() {
			<-network.ServerReady
			SyncCmd.Run(cmd, args)
		}()
		network.StartServer()
	},
}
