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
)

func call(node, method string, args interface{}, reply interface{}) error {
	log.Infof("call %s's %s, args: %+v", node, method, args)
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
		log.Warnln("call", ret.Message)
		return errors.New(ret.Message)
	}
	return nil
}

func CallBack(node, method string, args interface{}, reply interface{}) error {
	log.Debugln("Callback", method)
	return call(node, method, args, reply)
}

func CallSelf(method string, args interface{}, reply interface{}) error {
	log.Debugln("CallMySelf", method)
	return call("127.0.0.1:"+global.Port, method, args, reply)
}

// 左闭右开,target in [l,r)
func in(target, l, r int) bool {
	return l <= target && target < r
}

func CallInterGroup(method string, args interface{}, reply interface{}) (error, string) {
	log.Debugln("CallInterGroup", method)
	for _, node := range GetSortedNodes() {
		err := call(node.Address, method, args, reply)
		if err != nil {
			log.Warnln("CallInterGroup", node.Address, err)
			continue
		}
		return nil, node.Address
	}
	return errors.New("None of the nodes responded!"), ""
}

func CallInnerGroup(method string, args interface{}, reply interface{}) (error, string) {
	log.Debugln("CallInnerGroup", method)
	for _, node := range GetSortedNodes() {
		// 分组检查
		if node.GroupBase+node.GroupNumber > global.MaxGroupNum {
			if !in(global.GetGroup(), node.GroupBase, global.MaxGroupNum) &&
				!in(global.GetGroup(), 0, node.GroupBase+node.GroupNumber-global.MaxGroupNum) {
				continue
			}
		} else {
			if !in(global.GetGroup(), node.GroupBase, node.GroupBase+node.GroupNumber) {
				continue
			}
		}

		err := call(node.Address, method, args, reply)
		if err != nil {
			log.Warnln("CallInnerGroup", node.Address, err)
			continue
		}
		return nil, node.Address
	}
	return errors.New("None of the nodes responded!"), ""
}

func send(node Position, method string, args interface{}, reply interface{}) bool {
	// 分组检查
	if node.GroupBase+node.GroupNumber > global.MaxGroupNum {
		if !in(global.GetGroup(), node.GroupBase, global.MaxGroupNum) &&
			!in(global.GetGroup(), 0, node.GroupBase+node.GroupNumber-global.MaxGroupNum) {
			return false
		}
	} else {
		if !in(global.GetGroup(), node.GroupBase, node.GroupBase+node.GroupNumber) {
			return false
		}
	}

	err := call(node.Address, method, args, reply)
	if err != nil {
		log.Debugln("GossipCall", node.Address, method, err)
	} else {
		log.Debugln("GossipCall", node.Address, method, "success!")
	}
	return err == nil
}

func GossipCallSpecialGroup(method string, args interface{}, reply interface{}, targetGroup int) error {
	log.Debugln("GossipCall", "start", method, targetGroup)

	// 打乱数组
	nodes := sortedNodes
	for i := range nodes {
		rand.Seed(time.Now().UnixNano())
		j := rand.Int()%(len(nodes)-i) + i
		nodes[i], nodes[j] = nodes[j], nodes[i]
	}

	success := 0
	for _, node := range nodes {
		if send(node, method, args, reply) {
			success++
		}
		if success >= 3 {
			break
		}
	}

	if success == 0 {
		log.Warnln("GossipCall", "None of the nodes responded!")
		return errors.New("None of the nodes responded!")
	}

	return nil
}

func GossipCallInnerGroup(method string, args interface{}, reply interface{}) error {
	return GossipCallSpecialGroup(method, args, reply, global.GetGroup())
}

func GossipCallInterGroup(method string, args interface{}, reply interface{}) error {
	return GossipCallSpecialGroup(method, args, reply, -1)
}
