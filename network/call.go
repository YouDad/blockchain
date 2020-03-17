package network

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/YouDad/blockchain/log"
)

func call(node, method string, args interface{}, reply interface{}) error {
	log.Infoln("Request", method, node, args)
	b, err := json.Marshal(args)
	if err != nil {
		return err
	}
	resp, err := http.Post(fmt.Sprintf("http://%s/v1/%s?address=127.0.0.1:%s", node, method, Port),
		"application/json;charset=UTF-8", bytes.NewBuffer(b))
	if err != nil {
		return err
	}

	type SimpleJSONResult struct {
		Message string
		Data    interface{}
	}
	var ret SimpleJSONResult
	ret.Data = reply

	json.NewDecoder(resp.Body).Decode(&ret)
	if ret.Message != "" {
		log.Warnln(ret.Message)
	}
	return nil
}

func Callback(node, method string, args interface{}, reply interface{}) error {
	log.Infoln("Callback", method)
	return call(node, method, args, reply)
}

func Call(method string, args interface{}, reply interface{}) (error, string) {
	log.Infoln("Call", method)
	for _, node := range GetSortedNodes() {
		err := call(node.Address, method, args, reply)
		if err != nil {
			log.Warnln(node.Address, err)
			continue
		}
		return nil, node.Address
	}
	return errors.New("None of the nodes responded!"), ""
}

func CallMySelf(method string, args interface{}, reply interface{}) error {
	log.Infoln("CallMySelf", method)
	return call("127.0.0.1:"+Port, method, args, reply)
}

func random(min, max int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Int()%(max-min) + min
}

func GossipCall(method string, args interface{}, reply interface{}) error {
	log.Infoln("GossipCall", method)
	visit := make([]bool, len(sortedNodes))
	visited := 0
	success := 0

	send := func() bool {
		for {
			visitor := random(0, len(sortedNodes))
			if visit[visitor] {
				continue
			}
			visit[visitor] = true
			visited++
			err := call(sortedNodes[visitor].Address, method, args, reply)
			if err != nil {
				log.Infoln(sortedNodes[visitor].Address, err)
				continue
			}
			log.Infoln(sortedNodes[visitor].Address, "success!")
			return true
		}
	}

	for success < 3 && visited < len(sortedNodes) {
		if send() {
			success++
		}
	}

	if success == 0 {
		log.Warnln("None of the nodes responded!")
		return errors.New("None of the nodes responded!")
	}

	return nil
}
