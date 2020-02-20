package main

import (
	"log"

	"github.com/YouDad/blockchain/app/coin"
)

func main() {
	log.SetFlags(log.Lshortfile | log.Lmicroseconds | log.Ltime)
	log.SetPrefix("[info]: ")

	coin.Init()
	coin.Main()
}
