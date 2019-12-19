package main

import (
	"github.com/YouDad/blockchain/app/coin"
	"github.com/YouDad/blockchain/app/text"
)

func main() {
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
