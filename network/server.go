package network

import (
	"fmt"
	"sync"
	"time"

	"github.com/YouDad/blockchain/global"
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

func StartServer(sync func(group int) error) {
	// 周期性维持网络结构
	go func() {
		for {
			time.Sleep(5 * time.Second)

			// 获得节点列表
			GetKnownNodes()

			// 并行发送心跳包
			knownNodes := global.GetKnownNodes()
			ready := make(chan interface{}, 1)
			var nodeNumber int

			for nodeAddress := range knownNodes.Get() {
				nodeNumber++
				go func(address string) {
					start := time.Now().UnixNano()
					heartBeat(address)
					end := time.Now().UnixNano()
					knownNodes.UpdateNode(address, end-start)
					ready <- 0
				}(nodeAddress)
			}

			knownNodes.Release()

			for i := 0; i < nodeNumber; i++ {
				<-ready
			}
			UpdateSortedNodes()

			time.Sleep(10 * time.Second)
		}
	}()

	// 周期性同步节点
	go func() {
		for {
			time.Sleep(5 * time.Second)

			// 并行同步
			ready := make(chan interface{}, 1)

			for i := 0; i < global.GroupNum; i++ {
				go func(group int) {
					sync(group)
					ready <- 0
				}(i + global.GetGroup())
			}

			for i := 0; i < global.GroupNum; i++ {
				<-ready
			}

			time.Sleep(10 * time.Second)
		}
	}()

	// 过半秒设置信号
	go func() {
		time.Sleep(time.Millisecond * 500)
		ServerReady <- 0
	}()

	beego.Run(fmt.Sprintf("0.0.0.0:%s", global.Port))
}
