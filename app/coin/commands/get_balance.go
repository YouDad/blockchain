package commands

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/YouDad/blockchain/app/coin/core"
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
		bc := core.NewBlockchain()
		defer bc.Close()

		balance := 0
		UTXOs := bc.FindUTXO(getBalanceAddress)

		for _, out := range UTXOs {
			balance += out.Value
		}

		fmt.Printf("Balance of '%s': %d\n", getBalanceAddress, balance)
	},
}
