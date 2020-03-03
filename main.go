package main

import (
	cmd "github.com/YouDad/blockchain/commands"
	"github.com/YouDad/blockchain/log"
)

// Run parses command line arguments and processes commands
func main() {
	rootCmd := cmd.RootCmd
	rootCmd.AddCommand(
		cmd.GetBalanceCmd,
		cmd.CreateBlockchainCmd,
		cmd.SendCmd,
		cmd.GetVersionCmd,
		cmd.ListAddressCmd,
		cmd.CreateWalletCmd,
		cmd.StartNodeCmd,
	)

	if err := rootCmd.Execute(); err != nil {
		log.Errln(err)
		return
	}
}
