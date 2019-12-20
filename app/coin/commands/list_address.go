package commands

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"

	"github.com/YouDad/blockchain/app/coin/wallet"
)

var ListAddressCmd = &cobra.Command{
	Use:   "list_address",
	Short: "Lists all addresses from the wallet file",
	Run: func(cmd *cobra.Command, args []string) {
		wallets, err := wallet.NewWallets()
		if err != nil {
			log.Panic(err)
		}

		addresses := wallets.GetAddresses()

		for _, address := range addresses {
			fmt.Println(address)
		}
	},
}
