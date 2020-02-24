package network

import (
	"sort"
	"time"

	"github.com/YouDad/blockchain/log"
)

var (
	knownNodes  = make(map[string][6]int)
	sortedNodes PositionSlice
)

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

func addKnownNode(nodeAddress string) {
	_, ok := knownNodes[nodeAddress]
	if !ok {
		knownNodes[nodeAddress] = [6]int{0, 0, 0, 0, 0, 0}
	}
}

func updateKnownNode(node string, nano int64) {
	arr := knownNodes[node]
	arr[arr[5]] = int(nano / 1e9)
	arr[5] = (arr[5] + 1) % 5
	knownNodes[node] = arr
}

func updateSortedNodes() {
	sortedNodes = nil
	for address, times := range knownNodes {
		time := 0
		for i := 0; i < 5; i++ {
			time += times[i]
		}
		sortedNodes = append(sortedNodes, Position{
			Address:  address,
			Distance: time / 5,
		})
	}
	sort.Sort(sortedNodes)
}

func GetKnownNodes() error {
	log.NotImplement()
	return nil
}

func GetSortedNodes() PositionSlice {
	log.NotImplement()
	return nil
}

func knownNodeUpdating() {
	for {
		time.Sleep(40 * time.Second)
		GetKnownNodes()
		for nodeAddress, _ := range knownNodes {
			address := nodeAddress
			go func() {
				start := time.Now().UnixNano()
				BOOL := true
				call(address, "NET.HeartBeat", &BOOL, &BOOL)
				end := time.Now().UnixNano()
				updateKnownNode(address, end-start)
			}()
		}
		time.Sleep(20 * time.Second)
		updateSortedNodes()
		log.Infof("Sorted %+v\n", sortedNodes)
	}
}
