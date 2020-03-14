package commands

import (
	"github.com/YouDad/blockchain/api"
	"github.com/YouDad/blockchain/log"
	"github.com/YouDad/blockchain/network"
	"github.com/spf13/cobra"
)

var GetVersionCmd = &cobra.Command{
	Use:   "get_version",
	Short: "Print version information the blocks of the blockchain",
	Run: func(cmd *cobra.Command, args []string) {
		network.Register(Port)
		version, err := api.GetVersion()
		if err != nil {
			network.Register(Port)
			go network.StartServer()
			<-network.ServerReady
			version, err = api.GetVersion()
			log.Err(err)
		}
		log.Infof("Version :%d\n", version.Version)
		log.Infof("Height  :%d\n", version.Height)
		log.Infof("NowHash :%x\n", version.NowHash)
		log.Infof("RootHash:%x\n", version.RootHash)
	},
}
