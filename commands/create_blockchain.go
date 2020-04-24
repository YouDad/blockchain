package commands

import (
	"github.com/YouDad/blockchain/core"
	"github.com/YouDad/blockchain/global"
	"github.com/YouDad/blockchain/log"
	"github.com/YouDad/blockchain/wallet"
	"github.com/spf13/cobra"
)

func init() {
	CreateBlockchainCmd.Flags().StringVar(&global.Address, "address", "",
		"The Address to be presented with Genesis Block Award")
	CreateBlockchainCmd.MarkFlagRequired("address")
}

var CreateBlockchainCmd = &cobra.Command{
	Use:   "create_blockchain",
	Short: "Create a blockchain and send genesis block reward to ADDRESS",
	Run: func(cmd *cobra.Command, args []string) {
		if !wallet.ValidateAddress(global.Address) {
			log.Errln("Address is not valid", global.Address)
		}
		log.Err(core.CreateBlockchain(global.Address))
		log.Infoln("Done!")
	},
}
