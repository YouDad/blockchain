package commands

import (
	"github.com/YouDad/blockchain/api"
	"github.com/YouDad/blockchain/global"
	"github.com/YouDad/blockchain/log"
	"github.com/YouDad/blockchain/network"
	"github.com/spf13/cobra"
)

func init() {
	GetVersionCmd.Flags().StringVar(&global.Address, "address", "", "node's address")
	GetVersionCmd.MarkFlagRequired("address")
}

var GetVersionCmd = &cobra.Command{
	Use:   "get_version",
	Short: "Print version information the blocks of the blockchain",
	Run: func(cmd *cobra.Command, args []string) {
		network.Register()
		version, err := api.GetVersion()
		if err != nil {
			network.Register()
			go network.StartServer()
			<-network.ServerReady
			version, err = api.GetVersion()
			log.Err(err)
		}
		log.Infof("Version :%d\n", version.Version)
		log.Infof("Height  :%d\n", version.Height)
		log.Infof("NowHash :%s\n", version.NowHash)
		log.Infof("RootHash:%s\n", version.RootHash)
	},
}
