package commands

import (
	"time"

	"github.com/YouDad/blockchain/api"
	"github.com/YouDad/blockchain/core"
	"github.com/YouDad/blockchain/global"
	"github.com/YouDad/blockchain/log"
	"github.com/YouDad/blockchain/network"
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
		network.Register()

		group := global.GetGroup()
		bc := core.GetBlockchain(group)

		if bc.GetHeight() < 0 {
			genesis, err := api.GetGenesis(group)
			for err != nil {
				log.Warn(err)
				time.Sleep(5 * time.Second)
				log.Warn(network.GetKnownNodes())
				network.UpdateSortedNodes()
				genesis, err = api.GetGenesis(group)
			}
			bc.AddBlock(genesis)
			core.GetUTXOSet(group).Reindex()
		}

		genesis := bc.GetGenesis()
		lastest := bc.GetLastest()
		var height int32
		var address string
		var err error
		for {
			height, err, address = api.SendVersion(group, lastest.Height, genesis.Hash(), lastest.Hash())
			log.Warn(err)
			if err == nil {
				break
			}
			time.Sleep(5 * time.Second)
		}

		api.SyncBlocks(group, height, address)
	},
}
