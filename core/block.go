package core

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
	"time"
)

type Block struct {
	Timestamp     int64
	PrevBlockHash []byte
	Hash          []byte
	Nonce         int64
	Height        int
	App           CoinApp
}

func NewBlock(data CoinApp, prevBlockHash []byte, height int) *Block {
	block := &Block{
		Timestamp:     time.Now().Unix(),
		PrevBlockHash: prevBlockHash,
		Hash:          []byte{},
		App:           data,
		Height:        height,
	}

	pow := NewProofOfWork(block)
	nonce, hash := pow.Run()
	block.Hash = hash[:]
	block.Nonce = nonce
	return block
}

func (b *Block) Serialize() []byte {
	var result bytes.Buffer

	err := gob.NewEncoder(&result).Encode(b)
	if err != nil {
		log.Panic(err)
	}

	return result.Bytes()
}

func DeserializeBlock(d []byte) *Block {
	if d == nil {
		log.Println("Failed: input is nil")
		return nil
	}
	var block Block

	err := gob.NewDecoder(bytes.NewReader(d)).Decode(&block)
	if err != nil {
		log.Panic(err)
	}

	return &block
}

func (block *Block) String() string {
	ret := ""
	ret += fmt.Sprintf("Height: %d\n", block.Height)
	ret += fmt.Sprintf("Prev: %x\n", block.PrevBlockHash)
	ret += fmt.Sprintf("Hash: %x\n", block.Hash)
	ret += fmt.Sprintf("Txs : %s\n\n", block.App.ToString())
	return ret
}
