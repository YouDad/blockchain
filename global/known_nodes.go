package global

import (
	"sync"

	"github.com/YouDad/blockchain/types"
	// "github.com/YouDad/blockchain/log"
)

var instance *KnownNodes
var once sync.Once

func GetKnownNodes() *KnownNodes {
	once.Do(func() {
		instance = new(KnownNodes)
		instance.nodes = make(map[string]NetworkNode)
	})
	return instance
}

type NetworkNode struct {
	ReactTime types.Queue
	Timestamp int64
	Groups    []int
}

type KnownNodes struct {
	nodes map[string]NetworkNode
	mutex sync.Mutex
}

func (this *KnownNodes) Get() map[string]NetworkNode {
	// log.Debugln("KN Lock")
	this.mutex.Lock()
	return this.nodes
}

func (this *KnownNodes) Release() {
	this.mutex.Unlock()
	// log.Debugln("KN Unlock")
}

func (this *KnownNodes) AddNode(address string, time int64, group []int) {
	knownNodes := this.Get()
	defer this.Release()
	_, ok := knownNodes[address]
	if !ok {
		knownNodes[address] = NetworkNode{types.NewQueue(5), time, group}
	}
}

func (this *KnownNodes) UpdateNode(address string, nano int64) {
	knownNodes := this.Get()
	defer this.Release()
	node := knownNodes[address]
	if node.ReactTime.Len() == 5 {
		node.ReactTime.Pop()
	}
	node.ReactTime.Push(nano / 1e9)
	knownNodes[address] = node
}
