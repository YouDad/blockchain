package commands

import (
	"github.com/YouDad/blockchain/global"
	"github.com/YouDad/blockchain/log"
	"github.com/YouDad/blockchain/wallet"
	"github.com/spf13/cobra"
)

var specified int

func init() {
	CreateWalletCmd.Flags().IntVar(&specified, "specified", -1, "use specified group")
}

var CreateWalletCmd = &cobra.Command{
	Use:   "create_wallet",
	Short: "Generates a new key-pair and saves it into the wallet file",
	Run: func(cmd *cobra.Command, args []string) {
		ws, err := wallet.GetWallets()
		log.Err(err)
		w := wallet.NewWallet()
		global.Address = w.String()
		for specified != -1 && specified != global.GetGroup() {
			w = wallet.NewWallet()
			global.Address = w.String()
		}
		ws[w.String()] = w
		ws.SaveToFile()

		log.Infof("Your new address: %s\n", global.Address)
	},
}
