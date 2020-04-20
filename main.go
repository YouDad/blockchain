package main

import (
	cmd "github.com/YouDad/blockchain/commands"
	"github.com/YouDad/blockchain/log"
	_ "github.com/YouDad/blockchain/routers"
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
		cmd.MiningCmd,
		cmd.SyncCmd,
		cmd.AllCmd,
		cmd.SendTestCmd,
		cmd.PrintCmd,
	)

	if err := rootCmd.Execute(); err != nil {
		log.Errln(err)
		return
	}
}
