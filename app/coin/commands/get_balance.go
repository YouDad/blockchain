package commands

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/YouDad/blockchain/app/coin/core"
	"github.com/YouDad/blockchain/utils"
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
		pubKeyHash := utils.Base58Decode([]byte(getBalanceAddress))
		pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-4]
		UTXOs := bc.FindUTXO(pubKeyHash)

		for _, out := range UTXOs {
			balance += out.Value
		}

		fmt.Printf("Balance of '%s': %d\n", getBalanceAddress, balance)
	},
}
