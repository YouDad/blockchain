package commands

import (
	"fmt"

	"github.com/YouDad/blockchain/log"
	"github.com/YouDad/blockchain/store"
	"github.com/YouDad/blockchain/wallet"
	"github.com/spf13/cobra"
)

var (
	Port     string
	LogLevel int
)

func init() {
	RootCmd.PersistentFlags().StringVar(&Port, "port", "",
		"The port is service's portid at localhost network")
	RootCmd.MarkPersistentFlagRequired("port")
	RootCmd.PersistentFlags().IntVarP(&LogLevel, "verbose", "v", 2,
		"Verbose information 0~3")
}

var RootCmd = &cobra.Command{
	Use:   "blockchain",
	Short: "Blockchain coin Application",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		store.RegisterDatabase(fmt.Sprintf("blockchain%s.db", Port))
		wallet.Register(Port)
		if LogLevel < 0 {
			LogLevel = 0
		}
		log.Register(LogLevel, Port)
	},
}
