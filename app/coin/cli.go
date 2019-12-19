package coin

import (
	"log"

	cmd "github.com/YouDad/blockchain/app/coin/commands"
	"github.com/YouDad/blockchain/app/coin/core"
)

func Init() {
	core.Init()
}

// Run parses command line arguments and processes commands
func Main() {
	rootCmd := cmd.RootCmd
	rootCmd.AddCommand(
		cmd.GetBalanceCmd,
		cmd.CreateBlockchainCmd,
		cmd.SendCmd,
		cmd.PrintChainCmd,
	)

	if err := rootCmd.Execute(); err != nil {
		log.Panic(err)
	}
}
