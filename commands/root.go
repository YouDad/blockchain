package commands

import (
	"fmt"

	"github.com/YouDad/blockchain/global"
	"github.com/YouDad/blockchain/log"
	"github.com/YouDad/blockchain/wallet"
	"github.com/spf13/cobra"
)

var (
	logLevel uint
)

func init() {
	RootCmd.PersistentFlags().StringVar(&global.Port, "port", "",
		"The port is service's portid at localhost network")
	RootCmd.MarkPersistentFlagRequired("port")
	RootCmd.PersistentFlags().UintVarP(&logLevel, "verbose", "v", 2,
		"Verbose information 0~3")
}

var RootCmd = &cobra.Command{
	Use:   "blockchain",
	Short: "Blockchain coin Application",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		global.RegisterDatabase(fmt.Sprintf("blockchain%s.db", global.Port))
		wallet.Register()
		log.Register(logLevel, global.Port)
	},
}
