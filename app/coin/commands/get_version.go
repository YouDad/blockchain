package commands

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/YouDad/blockchain/log"
	"github.com/YouDad/blockchain/rpc"
)

var GetVersionCmd = &cobra.Command{
	Use:   "get_version",
	Short: "Print version information the blocks of the blockchain",
	Run: func(cmd *cobra.Command, args []string) {
		version, err := rpc.GetVersion(Port)
		if err != nil {
			go rpc.StartServer(Port, "")
			<-rpc.ServerReady
			version, err = rpc.GetVersion(Port)
			if err != nil {
				log.Panic(err)
			}
		}
		fmt.Printf("Version :%d\n", version.Version)
		fmt.Printf("Height  :%d\n", version.Height)
		fmt.Printf("RootHash:%x\n", version.RootHash)
	},
}
