package global

import (
	"sync"

	"github.com/YouDad/blockchain/log"
)

var instance *KnownNodes
var once sync.Once

func GetKnownNodes() *KnownNodes {
	once.Do(func() {
		instance = new(KnownNodes)
		instance.nodes = make(map[string][6]int)
	})
	return instance
}

type KnownNodes struct {
	nodes map[string][6]int
	mutex sync.Mutex
}

func (this *KnownNodes) Get() map[string][6]int {
	log.Debugln("KN Lock")
	this.mutex.Lock()
	return this.nodes
}

func (this *KnownNodes) Release() {
	this.mutex.Unlock()
	log.Debugln("KN Unlock")
}

func (this *KnownNodes) AddNode(address string) {
	knownNodes := this.Get()
	defer this.Release()
	_, ok := knownNodes[address]
	if !ok {
		knownNodes[address] = [6]int{0, 0, 0, 0, 0, 0}
	}
}

func (this *KnownNodes) UpdateNode(address string, nano int64) {
	knownNodes := this.Get()
	defer this.Release()
	arr := knownNodes[address]
	arr[arr[5]] = int(nano / 1e9)
	arr[5] = (arr[5] + 1) % 5
	knownNodes[address] = arr
}
