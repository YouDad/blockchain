package commands

import (
	"github.com/YouDad/blockchain/core"
	"github.com/YouDad/blockchain/log"
	"github.com/YouDad/blockchain/wallet"
	"github.com/spf13/cobra"
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
			log.Errln("Address is not valid", createBlockchainAddress)
		}
		core.CreateBlockchain(createBlockchainAddress)
		core.GetUTXOSet().Reindex()
		log.Infoln("Done!")
	},
}
