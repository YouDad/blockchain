package commands

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/YouDad/blockchain/app/coin/wallet"
)

var CreateWalletCmd = &cobra.Command{
	Use:   "create_wallet",
	Short: "Generates a new key-pair and saves it into the wallet file",
	Run: func(cmd *cobra.Command, args []string) {
		wallets, _ := wallet.NewWallets()
		address := wallets.CreateWallet()
		wallets.SaveToFile()

		fmt.Printf("Your new address: %s\n", address)
	},
}
