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
		wallets, err := wallet.NewWallets(Port)
		if err != nil {
			log.Errln(err)
		}

		addresses := wallets.GetAddresses()

		for _, address := range addresses {
			log.Infoln(address)
		}
	},
}
