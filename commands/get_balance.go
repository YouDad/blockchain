package commands

import (
	"github.com/YouDad/blockchain/api"
	"github.com/YouDad/blockchain/global"
	"github.com/YouDad/blockchain/log"
	"github.com/YouDad/blockchain/network"
	"github.com/spf13/cobra"
)

func init() {
	GetBalanceCmd.Flags().StringVar(&global.Address, "address", "", "The address of balance")
	GetBalanceCmd.MarkFlagRequired("address")
}

var GetBalanceCmd = &cobra.Command{
	Use:   "get_balance",
	Short: "Get balance of ADDRESS",
	Run: func(cmd *cobra.Command, args []string) {
		network.Register()
		balance, err := api.GetBalance(global.Address)
		if err != nil {
			network.Register()
			go network.StartServer()
			<-network.ServerReady
			balance, err = api.GetBalance(global.Address)
			log.Err(err)
		}
		log.Infof("Balance of '%s': %d\n", global.Address, balance)
	},
}
