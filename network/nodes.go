package network

import (
	"sort"
	"time"

	"github.com/YouDad/blockchain/global"
	"github.com/YouDad/blockchain/log"
)

var sortedNodes PositionSlice

type Position struct {
	Address     string
	Distance    int
	GroupBase   int
	GroupNumber int
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

func UpdateSortedNodes() {
	sortedNodes = nil
	knownNodes := global.GetKnownNodes()
	defer knownNodes.Release()
	for address, node := range knownNodes.Get() {
		time := 0

		for _, rt := range node.ReactTime.Get() {
			reactTime, ok := rt.(int)
			if ok {
				time += reactTime
			}
		}

		sortedNodes = append(sortedNodes, Position{
			Address:     address,
			Distance:    time / 5,
			GroupBase:   node.GroupBase,
			GroupNumber: node.GroupNumber,
		})
	}
	sort.Sort(sortedNodes)
}

func GetKnownNodes() error {
	knownNodes := []GetKnownNodesArgs{}
	myAddress := "127.0.0.1:" + global.Port
	err := getKnownNodes(myAddress, &knownNodes)
	if err == nil {
		for _, node := range knownNodes {
			global.GetKnownNodes().AddNode(node.Address, node.Timestamp, node.GroupBase, node.GroupNumber)
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
		time.Sleep(20 * time.Second)
		UpdateSortedNodes()
		log.Debugf("Sorted %+v\n", sortedNodes)
	}
}
