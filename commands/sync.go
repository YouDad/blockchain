package commands

import (
	"github.com/YouDad/blockchain/api"
	"github.com/YouDad/blockchain/global"
	"github.com/YouDad/blockchain/log"
	"github.com/spf13/cobra"
)

func init() {
	SyncCmd.Flags().IntVar(&global.GroupNum, "group", 1, "process group of number")
	SyncCmd.Flags().StringVar(&global.Address, "address", "", "address of node")
	SyncCmd.MarkFlagRequired("address")
}

var SyncCmd = &cobra.Command{
	Use:   "sync",
	Short: "sync information from other node",
	Run: func(cmd *cobra.Command, args []string) {
		log.Debugln("Syncing", global.Port)
		group := global.GetGroup()
		api.Sync(group)
	},
}
