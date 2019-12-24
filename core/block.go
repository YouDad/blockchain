package core

import (
	"bytes"
	"encoding/gob"
	"time"

	"github.com/YouDad/blockchain/app"
	"github.com/YouDad/blockchain/log"
)

type Block struct {
	Timestamp     int64
	PrevBlockHash []byte
	Hash          []byte
	Nonce         int64
	Height        int
	App           app.App
}

func NewBlock(data app.App, prevBlockHash []byte, height int) *Block {
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
	var block Block
	block.App = CoreConfig.GetAppdata()

	err := gob.NewDecoder(bytes.NewReader(d)).Decode(&block)
	if err != nil {
		log.Panic(err)
	}

	return &block
}
