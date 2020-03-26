package commands

import (
	"github.com/YouDad/blockchain/api"
	"github.com/YouDad/blockchain/core"
	"github.com/YouDad/blockchain/log"
	"github.com/YouDad/blockchain/network"
	"github.com/YouDad/blockchain/types"
	"github.com/YouDad/blockchain/wallet"
	"github.com/spf13/cobra"
)

var (
	sendFrom   string
	sendTo     string
	sendAmount int64
	sendMine   bool
)

func init() {
	SendCmd.Flags().StringVar(&sendFrom, "from", "", "Source wallet address")
	SendCmd.Flags().StringVar(&sendTo, "to", "", "Destination wallet address")
	SendCmd.Flags().Int64Var(&sendAmount, "amount", 0, "Amount to send")
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
			log.Errln("Sender address is not valid")
		}

		if !wallet.ValidateAddress(sendTo) {
			log.Errln("Recipient address is not valid")
		}

		network.Register()
		if sendMine {
			bc := core.GetBlockchain()
			set := core.GetUTXOSet()

			tx, err := set.NewUTXOTransaction(sendFrom, sendTo, sendAmount)
			log.Err(err)
			cbTx := core.NewCoinbaseTxn(sendFrom)
			txs := []*types.Transaction{cbTx, tx}

			newBlocks := bc.MineBlock(txs)
			bc.AddBlock(newBlocks)
			set.Update(newBlocks)
			return
		}
		err := api.SendCMD(sendFrom, sendTo, sendAmount)

		if err != nil {
			log.Warnln(err)
		} else {
			log.Infoln("Success!")
		}
	},
}
