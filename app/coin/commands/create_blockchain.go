package commands

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"

	"github.com/YouDad/blockchain/app/coin/coin_core"
	"github.com/YouDad/blockchain/app/coin/wallet"
)

var (
	createBlockchainAddress string
)

func init() {
	CreateBlockchainCmd.Flags().StringVar(&createBlockchainAddress, "address", "",
		"The Address to be presented with Genesis Block Award")
	CreateBlockchainCmd.MarkFlagRequired("address")
}

var CreateBlockchainCmd = &cobra.Command{
	Use:   "create_blockchain",
	Short: "Create a blockchain and send genesis block reward to ADDRESS",
	Run: func(cmd *cobra.Command, args []string) {
		if !wallet.ValidateAddress(createBlockchainAddress) {
			log.Panic("Address is not valid")
		}
		coin_core.CreateBlockchain(createBlockchainAddress)
		utxoSet := coin_core.NewUTXOSet()
		defer utxoSet.Close()

		utxoSet.Reindex()
		fmt.Println("Done!")
	},
}
