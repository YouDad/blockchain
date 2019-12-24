package coin

import (
	cmd "github.com/YouDad/blockchain/app/coin/commands"
	"github.com/YouDad/blockchain/app/coin/core"
	"github.com/YouDad/blockchain/log"
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
		cmd.GetVersionCmd,
		cmd.ListAddressCmd,
		cmd.CreateWalletCmd,
		cmd.StartNodeCmd,
	)

	if err := rootCmd.Execute(); err != nil {
		log.Panic(err)
	}
}
