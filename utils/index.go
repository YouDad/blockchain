package utils

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/gob"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"reflect"

	"github.com/YouDad/blockchain/log"
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
	ret, err := json.Marshal(arg)
	log.Err(err)
	// var result bytes.Buffer
	// log.Err(gob.NewEncoder(&result).Encode(arg))
	// ret := result.Bytes()
	return ret
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

func IntIndexOf(arr []int, target int) int {
	for index, value := range arr {
		if value == target {
			return index
		}
	}
	return -1
}
