package global

import (
	"strconv"

	"github.com/YouDad/blockchain/log"
)

var (
	GroupNum int
	Port     string
)

func GetGroup() int {
	port, err := strconv.Atoi(Port)
	log.Err(err)
	return port / 1000
}
