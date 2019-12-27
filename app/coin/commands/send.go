package commands

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"

	"github.com/YouDad/blockchain/app/coin/core"
	"github.com/YouDad/blockchain/app/coin/wallet"
	"github.com/YouDad/blockchain/rpc"
)

var (
	sendFrom   string
	sendTo     string
	sendAmount int
	sendMine   bool
)

func init() {
	SendCmd.Flags().StringVar(&sendFrom, "from", "", "Source wallet address")
	SendCmd.Flags().StringVar(&sendTo, "to", "", "Destination wallet address")
	SendCmd.Flags().IntVar(&sendAmount, "amount", -1, "Amount to send")
	SendCmd.Flags().BoolVar(&sendMine, "mine", false, "")
	SendCmd.MarkFlagRequired("from")
	SendCmd.MarkFlagRequired("to")
	SendCmd.MarkFlagRequired("amount")
}

var SendCmd = &cobra.Command{
	Use:   "send",
	Short: "Send AMOUNT of coins from FROM address to TO",
	Run: func(cmd *cobra.Command, args []string) {
		if !wallet.ValidateAddress(sendFrom) {
			log.Panic("ERROR: Sender address is not valid")
		}

		if !wallet.ValidateAddress(sendTo) {
			log.Panic("ERROR: Recipient address is not valid")
		}

		bc := core.NewBlockchain()
		defer bc.Close()
		utxoSet := core.NewUTXOSet()
		defer utxoSet.Close()

		tx := utxoSet.NewUTXOTransaction(sendFrom, sendTo, sendAmount)

		if sendMine {
			cbTx := core.NewCoinbaseTX(sendFrom, "")
			txs := []*core.Transaction{cbTx, tx}

			newBlocks := bc.MineBlock(txs)
			utxoSet.Update(newBlocks)
		} else {
			rpc.SendTransaction(tx)
		}
		fmt.Println("Success!")
	},
}
