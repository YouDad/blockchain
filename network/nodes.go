package network

import (
	"sort"
	"time"

	"github.com/YouDad/blockchain/global"
	"github.com/YouDad/blockchain/log"
)

var sortedNodes PositionSlice

type Position struct {
	Address  string
	Distance int
}

type PositionSlice []Position

func (p PositionSlice) Len() int {
	return len(p)
}

func (p PositionSlice) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func (p PositionSlice) Less(i, j int) bool {
	return p[i].Distance < p[j].Distance
}

func updateSortedNodes() {
	sortedNodes = nil
	knownNodes := global.GetKnownNodes()
	defer knownNodes.Release()
	for address, times := range knownNodes.Get() {
		time := 0
		for i := 0; i < 5; i++ {
			time += times[i]
		}
		if address == "127.0.0.1:9999" {
			time += 100
		}
		sortedNodes = append(sortedNodes, Position{
			Address:  address,
			Distance: time / 5,
		})
	}
	sort.Sort(sortedNodes)
}

func GetKnownNodes() error {
	knownNodeAddresses := []string{}
	myAddress := "127.0.0.1:" + Port
	err := getKnownNodes(myAddress, &knownNodeAddresses)
	if err == nil {
		for _, address := range knownNodeAddresses {
			global.GetKnownNodes().AddNode(address)
		}
	}
	return err
}

func GetSortedNodes() PositionSlice {
	return sortedNodes
}

func knownNodeUpdating() {
	for {
		time.Sleep(20 * time.Second)
		GetKnownNodes()
		knownNodes := global.GetKnownNodes()
		for nodeAddress, _ := range knownNodes.Get() {
			address := nodeAddress
			go func() {
				start := time.Now().UnixNano()
				heartBeat(address)
				end := time.Now().UnixNano()
				knownNodes.UpdateNode(address, end-start)
			}()
		}
		knownNodes.Release()
		time.Sleep(20 * time.Second)
		updateSortedNodes()
		log.Infof("Sorted %+v\n", sortedNodes)
	}
}
