package main

import (
	"log"

	cmd "github.com/YouDad/blockchain/commands"
)

// Run parses command line arguments and processes commands
func main() {
	log.SetFlags(log.Lshortfile | log.Lmicroseconds | log.Ltime)
	log.SetPrefix("[info]: ")

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
		log.Panic(err)
	}
}
