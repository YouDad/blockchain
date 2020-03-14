package api

import (
	"bytes"
	"sync"

	"github.com/YouDad/blockchain/conf"
	"github.com/YouDad/blockchain/core"
	"github.com/YouDad/blockchain/log"
)

var syncMutex sync.Mutex

func syncBlocks(newHeight int32, address string) {
	syncMutex.Lock()
	set := core.GetUTXOSet()
	lastest := set.GetLastest()
	originHash := lastest.Hash
	lastestHeight := lastest.Height

	if newHeight > lastestHeight {
		var l int32 = 0
		var r int32 = lastestHeight
		log.Traceln("l", l, "r", r)
		for r >= l {
			mid := l + (r-l)/2
			log.Traceln("mid", mid)
			block := core.BytesToBlock(set.SetTable(conf.BLOCKS).Get(mid))
			hash, err := CallbackGetHash(mid, address)
			if err != nil {
				log.Warn(err)
				return
			}

			log.Tracef("%x === %x\n", hash, block.Hash)
			if bytes.Compare(hash, block.Hash) == 0 {
				l = mid + 1
			} else {
				r = mid - 1
			}
		}

		log.Traceln("l", l, "r", r)
		lastestBytes := set.SetTable(conf.BLOCKS).Get(r)
		set.SetTable(conf.BLOCKS).Set("lastest", lastestBytes)
		lastest := core.BytesToBlock(lastestBytes)
		if lastest == nil {
			log.Errln("二分nil")
		}
		lastestHash := lastest.Hash
		blocks, err := CallbackGetBlocks(l, newHeight, lastestHash, address)
		if err != nil {
			log.Warn(err)
			lastestBytes = set.SetTable(conf.BLOCKS).Get(originHash)
			set.SetTable(conf.BLOCKS).Set("lastest", lastestBytes)
			return
		}
		// TODO utxoSet.Delete(block)

		for _, block := range blocks {
			if bytes.Compare(block.PrevHash, lastestHash) == 0 {
				set.AddBlock(block)
				lastestHash = block.Hash
			} else {
				break
			}
		}
		set.Reindex()
	}
	syncMutex.Unlock()
}
