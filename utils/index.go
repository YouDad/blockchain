package utils

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/gob"
	"io/ioutil"
	"net/http"
	"reflect"
	"sync"

	"github.com/YouDad/blockchain/log"
	jsoniter "github.com/json-iterator/go"
)

func GetExternIP() string {
	resp, err := http.Get("http://myexternalip.com/raw")
	if err != nil {
		return ""
	}
	defer resp.Body.Close()
	content, _ := ioutil.ReadAll(resp.Body)
	return string(content)
}

func InterfaceIsNil(i interface{}) bool {
	defer func() {
		recover()
	}()
	vi := reflect.ValueOf(i)
	return vi.IsNil()
}

func Encode(arg interface{}) []byte {
	ret, err := jsoniter.Marshal(arg)
	log.Err(err)
	return ret
}

var mutexDecode sync.Mutex

func Decode(b []byte, v interface{}) (err error) {
	mutexDecode.Lock()
	defer func() {
		if r := recover(); r != nil {
			log.Warnln(r)
			defer func() {
				log.SetCallerLevel(3)
				log.Errln(r)
			}()
			err = jsoniter.Unmarshal(b, v)
		}
	}()
	err = jsoniter.Unmarshal(b, v)
	mutexDecode.Unlock()
	return err
}

func SHA256(arg interface{}) []byte {
	hash := sha256.Sum256(Encode(arg))
	return hash[:]
}

func BaseTypeToBytes(num interface{}) []byte {
	switch number := num.(type) {
	case int:
		num = int32(number)
	}
	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.BigEndian, num)
	if err != nil {
		log.Errln(err)
	}
	return buff.Bytes()
}

func GetDecoder(b []byte) *gob.Decoder {
	reader := bytes.NewReader(b)
	return gob.NewDecoder(reader)
}

// 左闭右开,target in [l,r)
func in(target, l, r int) bool {
	return l <= target && target < r
}

func InGroup(group, base, number, max int) bool {
	// 分组检查
	if base+number > max {
		if !in(group, base, max) &&
			!in(group, 0, base+number-max) {
			return false
		}
	} else {
		if !in(group, base, base+number) {
			return false
		}
	}
	return true
}
