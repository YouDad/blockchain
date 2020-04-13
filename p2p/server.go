package p2p

import (
	"fmt"

	"github.com/YouDad/blockchain/log"
	"github.com/YouDad/blockchain/utils"
)

var (
	Port        string
	ServerReady = make(chan interface{}, 1)
)

func Register(port string) {
	Port = port
	addKnownNode("127.0.0.1:9999")
	externIP := utils.GetExternIP()
	addKnownNode(fmt.Sprintf("%s:%s", externIP, port))
	updateSortedNodes()
}

func StartServer(minerAddress string) {
	log.Errln("NotImplement")
	// utxoSet := core.NewUTXOSet()
	// go mining(minerAddress, utxoSet)
	// go knownNodeUpdating()
	// rpc.Register(&DB{utxoSet})
	// rpc.Register(&NET{utxoSet})
	// rpc.HandleHTTP()
	// l, err := net.Listen(protocol, fmt.Sprintf("0.0.0.0:%s", port))
	// if err != nil {
	//     log.Errln("listen error:", err)
	// }
	go func() { ServerReady <- 0 }()
	/* http.Serve(l, nil) */
}
