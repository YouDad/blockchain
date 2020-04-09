package commands

import (
	"github.com/YouDad/blockchain/global"
	"github.com/YouDad/blockchain/log"
	"github.com/YouDad/blockchain/wallet"
	"github.com/spf13/cobra"
)

var specified int

func init() {
	CreateWalletCmd.Flags().IntVar(&specified, "special_group", -1, "use specified group")
}

var CreateWalletCmd = &cobra.Command{
	Use:   "create_wallet",
	Short: "Generates a new key-pair and saves it into the wallet file",
	Run: func(cmd *cobra.Command, args []string) {
		wallets, _ := wallet.NewWallets()
		global.Address = wallets.CreateWallet()
		for specified != -1 && specified != global.GetGroup() {
			wallets, _ = wallet.NewWallets()
			global.Address = wallets.CreateWallet()
		}
		wallets.SaveToFile()

		log.Infof("Your new address: %s\n", global.Address)
	},
}
