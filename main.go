package main

import (
	"log"

	"github.com/YouDad/blockchain/coin_core"
	cmd "github.com/YouDad/blockchain/commands"
)

func init() {
	coin_core.Init()
}

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
