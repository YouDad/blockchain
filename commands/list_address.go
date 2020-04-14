package commands

import (
	"github.com/YouDad/blockchain/log"
	"github.com/YouDad/blockchain/wallet"
	"github.com/spf13/cobra"
)

var ListAddressCmd = &cobra.Command{
	Use:   "list_address",
	Short: "Lists all addresses from the wallet file",
	Run: func(cmd *cobra.Command, args []string) {
		wallets, err := wallet.GetWallets()
		log.Err(err)

		for address := range wallets {
			log.Infoln(address)
		}
	},
}
