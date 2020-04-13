package global

import (
	"github.com/YouDad/blockchain/utils"
)

var (
	GroupNum    int
	Port        string
	Address     string
	MaxGroupNum int
)

func GetGroup() int {
	pubKeyHash := utils.Base58Decode([]byte(Address))
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-4]
	return int(pubKeyHash[0]) % MaxGroupNum
}
