package rpc

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"

	coin_core "github.com/YouDad/blockchain/app/coin/core"
)

const (
	protocol = "tcp"
	version  = 1
)

var (
	ServerReady = make(chan interface{}, 1)
)

func StartServer(port, minerAddress string) {
	externIP := getExternIP()
	addKnownNode(fmt.Sprintf("%s:%s", externIP, port))
	addKnownNode("223.104.175.242:9999")
	updateSortedNodes()

	utxoSet := coin_core.NewUTXOSet()
	go mining(minerAddress, utxoSet)
	go knownNodeUpdating()
	rpc.Register(&DB{utxoSet})
	rpc.Register(&NET{utxoSet})
	rpc.HandleHTTP()
	l, err := net.Listen(protocol, fmt.Sprintf("0.0.0.0:%s", port))
	if err != nil {
		log.Fatal("listen error:", err)
	}
	go func() { ServerReady <- 0 }()
	http.Serve(l, nil)
}