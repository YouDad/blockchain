package core

import (
	"bytes"
	"encoding/gob"
	"log"
	"time"

	"github.com/YouDad/blockchain/app"
)

type Block struct {
	Timestamp     int64
	PrevBlockHash []byte
	Hash          []byte
	Nonce         int64
	App           app.App
}

func NewBlock(data app.App, prevBlockHash []byte) *Block {
	block := &Block{
		Timestamp:     time.Now().Unix(),
		PrevBlockHash: prevBlockHash,
		Hash:          []byte{},
		App:           data,
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
	block.App = coreConfig.GetAppdata()

	err := gob.NewDecoder(bytes.NewReader(d)).Decode(&block)
	if err != nil {
		log.Panic(err)
	}

	return &block
}
