package rpc

import (
	coin_core "github.com/YouDad/blockchain/app/coin/core"
)

func mining(address string, utxoSet *coin_core.UTXOSet) {
	if address == "" {
		return
	}

	for {
		txs := []*coin_core.Transaction{coin_core.NewCoinbaseTX(address, "")}
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
