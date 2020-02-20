package commands

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"

	"github.com/YouDad/blockchain/core"
	"github.com/YouDad/blockchain/wallet"
)

var ListAddressCmd = &cobra.Command{
	Use:   "list_address",
	Short: "Lists all addresses from the wallet file",
	Run: func(cmd *cobra.Command, args []string) {
		wallets, err := wallet.NewWallets(core.CoreConfig.WalletFile)
		if err != nil {
			log.Panic(err)
		}

		addresses := wallets.GetAddresses()

		for _, address := range addresses {
			fmt.Println(address)
		}
	},
}
