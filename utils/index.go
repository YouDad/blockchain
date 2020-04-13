package utils

import (
	"crypto/sha256"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"

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

type SHA256Arg interface {
	GobEncode() ([]byte, error)
	GobDecode(bytes []byte) error
}

func SHA256(arg SHA256Arg) types.HashValue {
	bytes, err := arg.GobEncode()
	if err != nil {
		log.Panic(err)
	}
	return sha256.Sum256(bytes)
}
