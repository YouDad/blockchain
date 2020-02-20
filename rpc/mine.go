package rpc

import (
	"github.com/YouDad/blockchain/core"
)

func mining(address string, utxoSet *core.UTXOSet) {
	if address == "" {
		return
	}

	for {
		txs := []*core.Transaction{core.NewCoinbaseTX(address, "")}
		nowHeight := utxoSet.GetBestHeight()
		nowSize := getMempoolSize()
		txs = append(txs, getTransactions()...)

		for {
			if nowHeight != utxoSet.GetBestHeight() {
				break
			}
			if nowSize != getMempoolSize() {
				break
			}

			newBlocks := utxoSet.MineBlock(txs)
			utxoSet.Update(newBlocks)
			SendBlock(newBlocks)
		}
	}
}
