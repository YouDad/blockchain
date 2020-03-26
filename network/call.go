package network

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/YouDad/blockchain/global"
	"github.com/YouDad/blockchain/log"
	"github.com/YouDad/blockchain/utils"
)

func call(node, method string, args interface{}, reply interface{}) error {
	log.Infoln("Request", method, node, args)
	b, err := json.Marshal(args)
	if err != nil {
		return err
	}
	resp, err := http.Post(fmt.Sprintf("http://%s/v1/%s?address=127.0.0.1:%s", node, method, global.Port),
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

func CallBack(node, method string, args interface{}, reply interface{}) error {
	log.Infoln("Callback", method)
	return call(node, method, args, reply)
}

func CallSelf(method string, args interface{}, reply interface{}) error {
	log.Infoln("CallMySelf", method)
	return call("127.0.0.1:"+global.Port, method, args, reply)
}

func CallInnerGroup(method string, args interface{}, reply interface{}) (error, string) {
	log.Infoln("Call", method)
	for _, node := range GetSortedNodes() {
		// TODO: use for to check Groups
		if utils.IntIndexOf(node.Groups, GetGroup()) == -1 {
			continue
		}
		err := call(node.Address, method, args, reply)
		if err != nil {
			log.Warnln(node.Address, err)
			continue
		}
		return nil, node.Address
	}
	return errors.New("None of the nodes responded!"), ""
}

func gossipCall(method string, args interface{}, reply interface{}, targetGroup int) error {
	log.Infoln("GossipCall", method)
	visit := make([]bool, len(sortedNodes))
	visited := 0
	success := 0

	send := func() bool {
		for len(sortedNodes) > visited {
			rand.Seed(time.Now().UnixNano())
			visitor := rand.Int() % len(sortedNodes)

			if visit[visitor] {
				continue
			}

			visit[visitor] = true
			visited++
			if targetGroup >= 0 && utils.IntIndexOf(sortedNodes[visitor].Groups, targetGroup) < 0 {
				continue
			}
			err := call(sortedNodes[visitor].Address, method, args, reply)
			if err != nil {
				log.Infoln(sortedNodes[visitor].Address, err)
				continue
			}
			log.Infoln(sortedNodes[visitor].Address, "success!")
			return true
		}
		return false
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

func GossipCallInnerGroup(method string, args interface{}, reply interface{}) error {
	return gossipCall(method, args, reply, GetGroup())
}

func GossipCallInterGroup(method string, args interface{}, reply interface{}) error {
	return gossipCall(method, args, reply, -1)
}
