package commands

import (
	"github.com/YouDad/blockchain/api"
	"github.com/YouDad/blockchain/log"
	"github.com/YouDad/blockchain/network"
	"github.com/spf13/cobra"
)

var (
	getBalanceAddress string
)

func init() {
	GetBalanceCmd.Flags().StringVar(&getBalanceAddress, "address", "", "The address of balance")
	GetBalanceCmd.MarkFlagRequired("address")
}

var GetBalanceCmd = &cobra.Command{
	Use:   "get_balance",
	Short: "Get balance of ADDRESS",
	Run: func(cmd *cobra.Command, args []string) {
		balance, err := api.GetBalance(getBalanceAddress)
		if err != nil {
			network.Register(Port)
			go network.StartServer()
			<-network.ServerReady
			balance, err = api.GetBalance(getBalanceAddress)
			if err != nil {
				log.Errln(err)
			}
		}
		log.Infof("Balance of '%s': %d\n", getBalanceAddress, balance)
	},
}
