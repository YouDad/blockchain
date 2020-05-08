package network

import (
	"math/rand"
	"time"

	"github.com/YouDad/blockchain/global"
)

var sortedNodes []Position

type Position struct {
	Address     string
	Distance    int
	GroupBase   int
	GroupNumber int
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

	rand.Seed(time.Now().UnixNano())
	for i := len(sortedNodes) - 1; i > 0; i-- {
		j := rand.Int() % i
		sortedNodes[i], sortedNodes[j] = sortedNodes[j], sortedNodes[i]
	}
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

func GetSortedNodes() []Position {
	return sortedNodes
}
