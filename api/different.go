package api

import (
	"bytes"

	"github.com/YouDad/blockchain/core"
	"github.com/YouDad/blockchain/global"
	"github.com/YouDad/blockchain/log"
)

func SyncBlocks(newHeight int32, address string) {
	lastestHeight := global.GetHeight()
	if newHeight <= lastestHeight {
		return
	}

	global.SyncMutex.Lock()
	defer global.SyncMutex.Unlock()
	log.Infoln("SyncBlock Start!")
	bc := core.GetBlockchain()
	lastest := bc.GetLastest()
	originHash := lastest.Hash

	if newHeight > lastestHeight {
		var l int32 = 0
		var r int32 = lastestHeight
		// log.Traceln("l", l, "r", r)
		for r >= l {
			mid := l + (r-l)/2
			// log.Traceln("mid", mid)
			block := core.BytesToBlock(bc.Get(mid))
			hash, err := CallbackGetHash(mid, address)
			if err != nil {
				log.Warn(err)
				return
			}

			// log.Tracef("%x === %x\n", hash, block.Hash)
			if bytes.Compare(hash, block.Hash) == 0 {
				l = mid + 1
			} else {
				r = mid - 1
			}
		}

		// log.Traceln("l", l, "r", r)
		lastestBytes := bc.Get(r)
		bc.Set("lastest", lastestBytes)
		global.SetHeight(lastestHeight)
		lastest := core.BytesToBlock(lastestBytes)
		if lastest == nil {
			log.Errln("二分nil")
		}

		lastestHash := lastest.Hash
		blocks, err := CallbackGetBlocks(l, newHeight, lastestHash, address)
		if err != nil {
			log.Warn(err)
			lastestBytes = bc.Get(originHash)
			bc.Set("lastest", lastestBytes)
			return
		}

		set := core.GetUTXOSet()
		for i := lastestHeight; i > r; i-- {
			set.Reverse(core.BytesToBlock(bc.Get(i)))
		}

		for _, block := range blocks {
			if bytes.Compare(block.PrevHash, lastestHash) == 0 {
				bc.AddBlock(block)
				lastestHash = block.Hash
			} else {
				break
			}
		}
	}
}
