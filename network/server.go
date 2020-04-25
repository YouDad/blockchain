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
	// 周期性维持网络结构
	go func() {
		for {
			time.Sleep(15 * time.Second)
			GetKnownNodes()
			knownNodes := global.GetKnownNodes()
			for nodeAddress := range knownNodes.Get() {
				address := nodeAddress
				go func() {
					start := time.Now().UnixNano()
					heartBeat(address)
					end := time.Now().UnixNano()
					knownNodes.UpdateNode(address, end-start)
				}()
			}
			knownNodes.Release()
			time.Sleep(15 * time.Second)
			UpdateSortedNodes()
			log.Debugf("Sorted %+v\n", sortedNodes)
		}
	}()

	// 过半秒设置信号
	go func() {
		time.Sleep(time.Millisecond * 500)
		ServerReady <- 0
	}()

	beego.Run(fmt.Sprintf("0.0.0.0:%s", global.Port))
}
