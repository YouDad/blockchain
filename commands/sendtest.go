package commands

import (
	"time"

	"github.com/YouDad/blockchain/api"
	"github.com/YouDad/blockchain/log"
	"github.com/YouDad/blockchain/network"
	"github.com/YouDad/blockchain/wallet"
	"github.com/spf13/cobra"
)

var (
	sendTestFrom string
	tps          int64
)

func init() {
	SendTestCmd.Flags().StringVar(&sendTestFrom, "from", "", "Source wallet address")
	SendTestCmd.Flags().Int64Var(&tps, "tps", 400, "send speed, transaction per second")
	SendTestCmd.MarkFlagRequired("from")
}

var SendTestCmd = &cobra.Command{
	Use:   "send_test",
	Short: "Send A lot of Txns for test, from FROM",
	Run: func(cmd *cobra.Command, args []string) {
		if !wallet.ValidateAddress(sendTestFrom) {
			log.Errln("Sender address is not valid")
		}

		network.Register()
		time.Sleep(time.Second)

		for {
			last := time.Now().UnixNano()
			sendTestTo := string(wallet.NewWallet().GetAddress())
			log.Infoln("SendTest", sendTestFrom, sendTestTo)
			err := api.SendCMD(sendTestFrom, sendTestTo, 1)

			if err != nil {
				log.Warnln("SendTest Warn?", err)
			} else {
				log.Infoln("SendTest Success!")
			}

			if 1e9/tps > time.Now().UnixNano()-last {
				time.Sleep(time.Duration(1e9/tps - (time.Now().UnixNano() - last)))
			}
		}
	},
}
