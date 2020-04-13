package commands

import (
	"github.com/spf13/cobra"
)

var (
	Port string
)

func init() {
	RootCmd.PersistentFlags().StringVar(&Port, "port", "",
		"The port is service's portid at localhost network")
	RootCmd.MarkPersistentFlagRequired("port")
}

var RootCmd = &cobra.Command{
	Use:   "blockchain",
	Short: "Blockchain coin Application",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// core.InitCore(core.Config{
		//     DatabaseFile: fmt.Sprintf("blockchain_%s.db", Port),
		//     WalletFile:   fmt.Sprintf("wallet_%s.dat", Port),
		// })
	},
}
