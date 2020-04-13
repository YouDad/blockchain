package p2p

import (
	"github.com/YouDad/blockchain/log"
)

type Position struct {
	Address  string
	Distance int
}

type PositionSlice []Position

func addKnownNode(nodeAddress string) {
	log.Errln("NotImplement")

}

func updateSortedNodes() {
	log.Errln("NotImplement")

}

func GetKnownNodes() error {
	log.Errln("NotImplement")

	return nil
}

func GetSortedNodes() PositionSlice {
	log.Errln("NotImplement")
	return nil
}
