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
	MiningCmd.Flags().IntVar(&global.GroupNum, "group", 1, "process group of number")
}

var MiningCmd = &cobra.Command{
	Use:   "mining",
	Short: "Start a mining node with ID specified in port.",
	Run: func(cmd *cobra.Command, args []string) {
		log.Infoln("Starting node", global.Port)
		network.Register()
		core.Register(speed)
		go func() {
			<-network.ServerReady
			if !wallet.ValidateAddress(nodeAddress) {
				log.Errln("Address is not valid")
			}

			group := global.GetGroup()
			for {
				var txns [][]*types.Transaction
				for i := 0; i < global.GroupNum; i++ {
					txns = append(txns, []*types.Transaction{core.NewCoinbaseTxn(nodeAddress)})
					txns[i] = append(txns[i], global.GetMempool().GetTxns(group+i)...)
				}

				for {
					log.Infoln("core.MineBlocks", group, global.GroupNum, "{{{{{{{{")
					newBlocks := core.MineBlocks(txns, group, global.GroupNum)
					if newBlocks == nil {
						break
					}
					for _, newBlock := range newBlocks {
						api.GossipBlock(newBlock)
					}
					log.Infoln("core.MineBlocks", group, global.GroupNum, "}}}}}}}}")
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
