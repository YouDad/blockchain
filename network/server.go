package network

import (
	"fmt"
	"time"

	"github.com/YouDad/blockchain/global"
	"github.com/YouDad/blockchain/log"
	"github.com/astaxie/beego"
)

var (
	Port        string
	protocol    = "tcp"
	ServerReady = make(chan interface{}, 1)
)

func Register(port string) {
	if Port == "" {
		Port = port
		global.GetKnownNodes().AddNode("127.0.0.1:9999")
		updateSortedNodes()
	}
}

func StartServer() {
	address := fmt.Sprintf("0.0.0.0:%s", Port)
	go knownNodeUpdating()
	log.Infoln("Server Listen", address)
	go func() {
		time.Sleep(time.Millisecond * 500)
		ServerReady <- 0
	}()
	beego.Run(address)
}
