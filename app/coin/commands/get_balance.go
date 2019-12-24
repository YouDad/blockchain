package commands

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/YouDad/blockchain/log"
	"github.com/YouDad/blockchain/rpc"
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
		balance, err := rpc.GetBalance(Port, getBalanceAddress)
		if err != nil {
			go rpc.StartServer(Port, "")
			<-rpc.ServerReady
			balance, err = rpc.GetBalance(Port, getBalanceAddress)
			if err != nil {
				log.Panic(err)
			}
		}
		fmt.Printf("Balance of '%s': %d\n", getBalanceAddress, balance)
	},
}
