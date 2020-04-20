package network

import (
	"fmt"
	"sync"
	"time"

	"github.com/YouDad/blockchain/global"
	"github.com/YouDad/blockchain/log"
	"github.com/astaxie/beego"
)

var (
	protocol     = "tcp"
	ServerReady  = make(chan interface{}, 1)
	onceRegister sync.Once
)

func Register() {
	onceRegister.Do(func() {
		global.GetKnownNodes().AddNode("127.0.0.1:1111", 0, 0, global.GroupNum)
		UpdateSortedNodes()
	})
}

func StartServer() {
	address := fmt.Sprintf("0.0.0.0:%s", global.Port)
	go knownNodeUpdating()
	log.Infoln("Server Listen", address)
	go func() {
		time.Sleep(time.Millisecond * 500)
		ServerReady <- 0
	}()
	beego.Run(address)
}
