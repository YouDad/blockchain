package commands

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/YouDad/blockchain/app/coin/core"
)

var PrintChainCmd = &cobra.Command{
	Use:   "print_chain",
	Short: "Print all the blocks of the blockchain",
	Run: func(cmd *cobra.Command, args []string) {
		bc := core.NewBlockchain()
		defer bc.Close()

		iter := bc.Begin()

		for {
			block := iter.Next()
			if block == nil {
				break
			}

			pow := core.NewProofOfWork(block)
			fmt.Printf("Prev: %x\n", block.PrevBlockHash)
			fmt.Printf("Hash: %x\n", block.Hash)
			fmt.Printf("PoW : %s\n", strconv.FormatBool(pow.Validate()))
			fmt.Printf("Txs : %s\n\n", block.App.ToString())
		}
	},
}
