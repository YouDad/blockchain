package commands

import (
	"time"

	"github.com/YouDad/blockchain/api"
	"github.com/YouDad/blockchain/global"
	"github.com/YouDad/blockchain/log"
	"github.com/YouDad/blockchain/network"
	"github.com/YouDad/blockchain/wallet"
	"github.com/spf13/cobra"
)

var tps int64
var wait int64

func init() {
	SendTestCmd.Flags().StringVar(&global.Address, "from", "", "Source wallet address")
	SendTestCmd.Flags().Int64Var(&tps, "tps", 10, "send speed, transaction per second")
	SendTestCmd.Flags().Int64Var(&wait, "wait", 90, "the time before sendloop")
	SendTestCmd.MarkFlagRequired("from")
}

var SendTestCmd = &cobra.Command{
	Use:   "send_test",
	Short: "Send A lot of Txns for test, from FROM",
	Run: func(cmd *cobra.Command, args []string) {
		if !wallet.ValidateAddress(global.Address) {
			log.Errln("Sender address is not valid")
		}

		network.Register()
		go network.StartServer()
		time.Sleep(time.Duration(wait) * time.Second)

		for {
			for global.GetMempool(global.GetGroup()).GetMempoolSize() >= 7*int(tps) {
				time.Sleep(time.Second)
			}

			sendTestTo := string(wallet.NewWallet().GetAddress())
			log.Infoln("SendTest", global.GetMempool(global.GetGroup()).GetMempoolSize(),
				global.Address, sendTestTo)
			err := api.SendCMD(global.Address, sendTestTo, 1)

			if err != nil {
				log.Warnln("SendTest Warn?", err)
			} else {
				log.Infoln("SendTest Success!")
			}
		}
	},
}
