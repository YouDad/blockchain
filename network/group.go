package network

import (
	"strconv"

	"github.com/YouDad/blockchain/global"
	"github.com/YouDad/blockchain/log"
)

func GetGroup() int {
	port, err := strconv.Atoi(global.Port)
	log.Err(err)
	return port / 1000
}
