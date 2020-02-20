package commands

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/YouDad/blockchain/core"
	"github.com/YouDad/blockchain/wallet"
)

var CreateWalletCmd = &cobra.Command{
	Use:   "create_wallet",
	Short: "Generates a new key-pair and saves it into the wallet file",
	Run: func(cmd *cobra.Command, args []string) {
		wallets, _ := wallet.NewWallets(core.CoreConfig.WalletFile)
		address := wallets.CreateWallet()
		wallets.SaveToFile(core.CoreConfig.WalletFile)

		fmt.Printf("Your new address: %s\n", address)
	},
}
