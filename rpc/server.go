package rpc

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"

	"github.com/YouDad/blockchain/app/coin/coin_core"
)

const (
	protocol = "tcp"
	version  = 1
)

var (
	ServerReady = make(chan interface{}, 1)
)

func Init(port string) {
	addKnownNode("39.107.64.93:9999")
	externIP := getExternIP()
	addKnownNode(fmt.Sprintf("%s:%s", externIP, port))
	updateSortedNodes()
}

func StartServer(port, minerAddress string) {
	Port = port
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
