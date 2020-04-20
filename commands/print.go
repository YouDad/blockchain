package commands

import (
	"github.com/YouDad/blockchain/core"
	"github.com/YouDad/blockchain/log"
	"github.com/spf13/cobra"
)

var group int

func init() {
	PrintCmd.Flags().IntVar(&group, "group", 0, "want to print db's group")
	PrintCmd.MarkFlagRequired("group")
}

var PrintCmd = &cobra.Command{
	Use:   "print",
	Short: "print the db",
	Run: func(cmd *cobra.Command, args []string) {
		bc := core.GetBlockchain(group)
		height := bc.GetHeight()
		var i int32
		for i = 0; i < height; i++ {
			block := bc.GetBlockByHeight(i)
			log.Infoln(block.Height, len(block.Txns))
		}
	},
}
