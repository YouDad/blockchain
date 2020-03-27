package api

import (
	"bytes"

	"github.com/YouDad/blockchain/core"
	"github.com/YouDad/blockchain/global"
	"github.com/YouDad/blockchain/log"
)

func SyncBlocks(group int, newHeight int32, address string) {
	lastestHeight := core.GetBlockchain().GetHeight(group)
	if newHeight <= lastestHeight {
		return
	}

	global.SyncMutex.Lock()
	defer global.SyncMutex.Unlock()
	log.Infoln("SyncBlock Start!")
	bc := core.GetBlockchain()
	lastest := bc.GetLastest(group)
	originHash := lastest.Hash()

	if newHeight > lastestHeight {
		var l int32 = 0
		var r int32 = lastestHeight
		// log.Traceln("l", l, "r", r)
		for r >= l {
			mid := l + (r-l)/2
			// log.Traceln("mid", mid)
			block := core.BytesToBlock(bc.Get(group, mid))
			hash, err := CallbackGetHash(group, mid, address)
			if err != nil {
				log.Warn(err)
				return
			}

			// log.Tracef("%s === %s\n", hash, block.Hash)
			if bytes.Compare(hash, block.Hash()) == 0 {
				l = mid + 1
			} else {
				r = mid - 1
			}
		}

		// log.Traceln("l", l, "r", r)
		lastestBytes := bc.Get(group, r)
		bc.Set(group, "lastest", lastestBytes)
		lastest := core.BytesToBlock(lastestBytes)
		if lastest == nil {
			log.Errln("二分nil")
		}

		lastestHash := lastest.Hash()
		blocks, err := CallbackGetBlocks(group, l, newHeight, lastestHash, address)
		if err != nil {
			log.Warn(err)
			lastestBytes = bc.Get(group, originHash)
			bc.Set(group, "lastest", lastestBytes)
			return
		}

		set := core.GetUTXOSet()
		for i := lastestHeight; i > r; i-- {
			set.Reverse(group, core.BytesToBlock(bc.Get(group, i)))
		}

		for _, block := range blocks {
			if bytes.Compare(block.PrevHash, lastestHash) == 0 {
				bc.AddBlock(group, block)
				set.Update(group, block)
				lastestHash = block.Hash()
			} else {
				break
			}
		}
	}
}
