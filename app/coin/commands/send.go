package commands

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/YouDad/blockchain/app/coin/core"
)

var (
	sendFrom   string
	sendTo     string
	sendAmount int
)

func init() {
	SendCmd.Flags().StringVar(&sendFrom, "from", "", "Source wallet address")
	SendCmd.Flags().StringVar(&sendTo, "to", "", "Destination wallet address")
	SendCmd.Flags().IntVar(&sendAmount, "amount", -1, "Amount to send")
	SendCmd.MarkFlagRequired("from")
	SendCmd.MarkFlagRequired("to")
	SendCmd.MarkFlagRequired("amount")
}

var SendCmd = &cobra.Command{
	Use:   "send",
	Short: "Send AMOUNT of coins from FROM address to TO",
	Run: func(cmd *cobra.Command, args []string) {
		bc := core.NewBlockchain()
		defer bc.Close()

		tx := bc.NewUTXOTransaction(sendFrom, sendTo, sendAmount)
		bc.AddBlock(core.GetCoinApp([]*core.Transaction{tx}))
		fmt.Println("Success!")
		fmt.Println("Done!")
	},
}
