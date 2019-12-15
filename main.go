package main

import (
	"fmt"
	"github.com/YouDad/blockchain/core"
)

func main() {
	bc := core.NewBlockchain()
	bc.AddBlock("1")
	bc.AddBlock("2")

	for _, block := range bc.Blocks {
		fmt.Printf("Prev: %x\n", block.PrevBlockHash)
		fmt.Printf("Hash: %x\n", block.Hash)
		fmt.Printf("Data: %s\n", block.Data)
		fmt.Println()
	}
}
