package rpc

import (
	"fmt"
	"log"
	"sort"
	"time"
)

type position struct {
	Address  string
	Distance int
}
type positionSlice []position

func (p positionSlice) Len() int {
	return len(p)
}
func (p positionSlice) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}
func (p positionSlice) Less(i, j int) bool {
	return p[i].Distance < p[j].Distance
}

var (
	knownNodes  = make(map[string][6]int)
	sortedNodes positionSlice
)

func updateSortedNodes() {
	sortedNodes = nil
	for address, times := range knownNodes {
		time := times[0]
		time += times[1]
		time += times[2]
		time += times[3]
		time += times[4]
		sortedNodes = append(sortedNodes, position{
			Address:  address,
			Distance: time / 5,
		})
	}
	sort.Sort(sortedNodes)
}

func addKnownNode(node string) {
	knownNodes[node] = [6]int{0, 0, 0, 0, 0, 0}
}

func updateKnownNode(node string, nano int64) {
	arr := knownNodes[node]
	arr[arr[5]] = int(nano / 1e9)
	arr[5] = (arr[5] + 1) % 5
	knownNodes[node] = arr
}

func delKnownNode(node string) {
	delete(knownNodes, node)
}

func knownNodeUpdating() {
	for {
		GetKnownNodes()
		for nodeAddress, _ := range knownNodes {
			address := nodeAddress
			go func() {
				start := time.Now().UnixNano()
				HeartBeat(address)
				end := time.Now().UnixNano()
				updateKnownNode(address, end-start)
			}()
		}
		time.Sleep(20 * time.Second)
		updateSortedNodes()
		time.Sleep(40 * time.Second)
	}
}

func GossipCall(method string, args interface{}, reply interface{}) {
	logger := fmt.Sprintf("GossipCall %s", method)
	log.Println(logger)
	len := len(sortedNodes)
	half := len / 2
	visit := make([]bool, len)

	canVisit := func(min, max int) bool {
		for _, v := range visit[min:max] {
			if !v {
				return true
			}
		}
		return false
	}

	send := func(min, max int) bool {
		for {
			if !canVisit(min, max) {
				return false
			}

			visitor := random(min, max)
			if visit[visitor] {
				continue
			}
			visit[visitor] = true
			_args := args
			_reply := reply
			err := call(sortedNodes[visitor].Address, method, _args, _reply)
			if err != nil {
				log.Println(logger, sortedNodes[visitor].Address, err)
				continue
			}
			log.Println(logger, sortedNodes[visitor].Address, "success!")
			return true
		}
	}

	if !send(0, half) {
		send(half+1, len)
	}
	if !send(0, half) {
		send(half+1, len)
	}
	send(half+1, len)
}
