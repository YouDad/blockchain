package rpc

import (
	"fmt"
	"net"
	"net/http"
	"net/rpc"

	coin_core "github.com/YouDad/blockchain/app/coin/core"
	"github.com/YouDad/blockchain/log"
)

const (
	protocol = "tcp"
	version  = 1
)

var (
	miningAddress string
	ServerReady   = make(chan interface{}, 1)
)

func StartServer(port, minerAddress string) {
	miningAddress = minerAddress
	utxo_set := coin_core.NewUTXOSet()
	rpc.Register(&DB{utxo_set})
	rpc.Register(&NET{utxo_set})
	rpc.HandleHTTP()
	l, err := net.Listen(protocol, fmt.Sprintf("0.0.0.0:%s", port))
	if err != nil {
		log.Fatal("listen error:", err)
	}
	go func() { ServerReady <- 0 }()
	http.Serve(l, nil)
}
