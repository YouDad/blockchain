package network

import (
	"fmt"
	"net"
	"net/http"
	"net/rpc"

	"github.com/YouDad/blockchain/log"
)

var (
	Port        string
	protocol    = "tcp"
	ServerReady = make(chan interface{}, 1)
)

func Register(port string) {
	if Port == "" {
		Port = port
		addKnownNode("127.0.0.1:9999")
		updateSortedNodes()
	}
}

func StartServer() {
	address := fmt.Sprintf("0.0.0.0:%s", Port)
	go knownNodeUpdating()
	rpc.HandleHTTP()
	l, err := net.Listen(protocol, address)
	log.Err(err)
	go func() { ServerReady <- 0 }()
	log.Infoln("Server Listen", address)
	http.Serve(l, nil)
}
