package network

import (
	"fmt"
	"net"
	"net/http"
	"net/rpc"

	"github.com/YouDad/blockchain/log"
	"github.com/YouDad/blockchain/utils"
)

var (
	Port        string
	protocol    = "tcp"
	ServerReady = make(chan interface{}, 1)
)

func Register(port string) {
	Port = port
	addKnownNode("127.0.0.1:9999")
	externIP := utils.GetExternIP()
	addKnownNode(fmt.Sprintf("%s:%s", externIP, port))
	addKnownNode(fmt.Sprintf("127.0.0.1:%s", port))
	updateSortedNodes()
}

func StartServer() {
	go knownNodeUpdating()
	rpc.HandleHTTP()
	l, err := net.Listen(protocol, fmt.Sprintf("0.0.0.0:%s", Port))
	log.Err(err)
	go func() { ServerReady <- 0 }()
	log.Infoln("Server Listen 0.0.0.0:" + Port)
	http.Serve(l, nil)
}
