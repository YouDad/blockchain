package commands

import (
	"fmt"

	"github.com/YouDad/blockchain/app/coin/core"
	"github.com/spf13/cobra"
)

var ReindexUTXOCmd = &cobra.Command{
	Use:   "reindex_utxo",
	Short: "Rebuilds the UTXO set",
	Run: func(cmd *cobra.Command, args []string) {
		UTXOSet := core.NewUTXOSet()
		UTXOSet.Reindex()

		count := UTXOSet.CountTransactions()
		fmt.Printf("Done! There are %d transactions in the UTXO set.\n", count)
	},
}
