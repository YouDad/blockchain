package commands

import (
	"github.com/YouDad/blockchain/log"
	"github.com/YouDad/blockchain/wallet"
	"github.com/spf13/cobra"
)

var CreateWalletCmd = &cobra.Command{
	Use:   "create_wallet",
	Short: "Generates a new key-pair and saves it into the wallet file",
	Run: func(cmd *cobra.Command, args []string) {
		wallets, _ := wallet.NewWallets()
		address := wallets.CreateWallet()
		wallets.SaveToFile(Port)

		log.Infof("Your new address: %s\n", address)
	},
}
