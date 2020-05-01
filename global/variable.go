package global

import (
	"github.com/YouDad/blockchain/types"
	"github.com/YouDad/blockchain/utils"
)

var (
	GroupNum    int
	Port        string
	Address     string
	MaxGroupNum int
)

// 返回默认组
func GetGroup() int {
	return GetGroupByAddress(Address)
}

// 返回地址对应的组
func GetGroupByAddress(address string) int {
	pubKeyHash := utils.Base58Decode([]byte(address))
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-4]
	return int(pubKeyHash[0]) % MaxGroupNum
}

// 返回公钥哈希对应的组
func GetGroupByPubKeyHash(pubKeyHash types.HashValue) int {
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-4]
	return int(pubKeyHash[0]) % MaxGroupNum
}
