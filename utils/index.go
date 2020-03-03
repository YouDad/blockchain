package utils

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/gob"
	"io/ioutil"
	"net/http"
	"reflect"

	"github.com/YouDad/blockchain/log"
	"github.com/YouDad/blockchain/types"
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
	var result bytes.Buffer
	log.Err(gob.NewEncoder(&result).Encode(arg))
	ret := result.Bytes()
	// log.SetCallerLevel(2)
	// log.Debugf("Encode %x", ret)
	// log.SetCallerLevel(0)
	return ret
}

func SHA256(arg interface{}) types.HashValue {
	hash := sha256.Sum256(Encode(arg))
	return hash[:]
}

func IntToBytes(num int64) []byte {
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
