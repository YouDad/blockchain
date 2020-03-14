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

func Call(method string, args interface{}, reply interface{}) error {
	log.Infoln("Call", method)
	for _, node := range GetSortedNodes() {
		err := call(node.Address, method, args, reply)
		if err != nil {
			log.Warnln(node.Address, err)
			continue
		}
		return nil
	}
	return errors.New("None of the nodes responded!")
}

func CallMySelf(method string, args interface{}, reply interface{}) error {
	log.Infoln("CallMySelf", method)
	return call("127.0.0.1:"+Port, method, args, reply)
}

func random(min, max int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Int()%(max-min) + min
}

func GossipCall(method string, args interface{}, reply interface{}) {
	log.Infoln("GossipCall", method)
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
			err := call(sortedNodes[visitor].Address, method, args, reply)
			if err != nil {
				log.Infoln(sortedNodes[visitor].Address, err)
				continue
			}
			log.Infoln(sortedNodes[visitor].Address, "success!")
			return true
		}
	}

	visited := 0

	if !send(0, half) {
		if send(half, len) {
			visited++
		}
	} else {
		visited++
	}

	if !send(0, half) {
		if send(half, len) {
			visited++
		}
	} else {
		visited++
	}

	if send(half, len) {
		visited++
	}

	if visited == 0 {
		log.Warnln("None of the nodes responded!")
	}
}
