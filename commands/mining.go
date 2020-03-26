package commands

import (
	"github.com/spf13/cobra"

	"github.com/YouDad/blockchain/api"
	"github.com/YouDad/blockchain/core"
	"github.com/YouDad/blockchain/global"
	"github.com/YouDad/blockchain/log"
	"github.com/YouDad/blockchain/network"
	"github.com/YouDad/blockchain/types"
	"github.com/YouDad/blockchain/wallet"
)

var (
	nodeAddress string
	speed       uint
)

func init() {
	MiningCmd.Flags().StringVar(&nodeAddress, "address", "", "node's coin address")
	MiningCmd.Flags().UintVar(&speed, "speed", 1, "mining speed: 0~100, 100 is 100% pc speed")
}

var MiningCmd = &cobra.Command{
	Use:   "mining",
	Short: "Start a mining node with ID specified in port.",
	Run: func(cmd *cobra.Command, args []string) {
		log.Infoln("Starting node", global.Port)
		network.Register()
		core.Register(speed)
		go func() {
			if !wallet.ValidateAddress(nodeAddress) {
				log.Errln("Address is not valid")
			}

			bc := core.GetBlockchain()
			for {
				txs := []*types.Transaction{core.NewCoinbaseTxn(nodeAddress)}
				set := core.GetUTXOSet()
				txs = append(txs, global.GetMempool().GetTxns()...)

				for {
					newBlocks := bc.MineBlock(txs)
					if newBlocks == nil {
						break
					}
					set.Update(newBlocks)
					api.GossipBlock(newBlocks)
				}
			}
		}()
		go func() {
			<-network.ServerReady
			SyncCmd.Run(cmd, args)
		}()
		network.StartServer()
	},
}
