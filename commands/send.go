package commands

import (
	"github.com/YouDad/blockchain/api"
	"github.com/YouDad/blockchain/core"
	"github.com/YouDad/blockchain/global"
	"github.com/YouDad/blockchain/log"
	"github.com/YouDad/blockchain/network"
	"github.com/YouDad/blockchain/types"
	"github.com/YouDad/blockchain/wallet"
	"github.com/spf13/cobra"
)

var (
	sendTo     string
	sendAmount int64
	sendMine   bool
)

func init() {
	SendCmd.Flags().StringVar(&global.Address, "from", "", "Source wallet address")
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
		if !wallet.ValidateAddress(global.Address) {
			log.Errln("Sender address is not valid")
		}

		if !wallet.ValidateAddress(sendTo) {
			log.Errln("Recipient address is not valid")
		}

		network.Register()
		if sendMine {
			bc := core.GetBlockchain(global.GetGroup())
			set := core.GetUTXOSet(global.GetGroup())

			tx, err := set.CreateTransaction(global.Address, sendTo, sendAmount)
			log.Err(err)
			cbTx := core.NewCoinbaseTxn(global.Address)
			txs := []*types.Transaction{cbTx, tx}

			newBlocks := core.MineBlocks([][]*types.Transaction{txs}, global.GetGroup(), 1)
			bc.AddBlock(newBlocks[0])
			set.Update(newBlocks[0])
			return
		}
		err := api.SendCMD(global.Address, sendTo, sendAmount)

		if err != nil {
			log.Warnln(err)
		} else {
			log.Infoln("Success!")
		}
	},
}
