package main

import (
	"log"

	"github.com/YouDad/blockchain/app/coin"
	"github.com/YouDad/blockchain/app/text"
)

func main() {
	log.SetFlags(log.Lshortfile | log.Lmicroseconds | log.Ltime)
	log.SetPrefix("[info]: ")

	app := 1
	switch app {
	case 0:
		text.Init()
		text.Main()
	case 1:
		coin.Init()
		coin.Main()
	default:
	}
}
