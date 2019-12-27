package rpc

import (
	"io/ioutil"
	"math/rand"
	"net/http"
	"time"
)

func getExternIP() string {
	resp, err := http.Get("http://myexternalip.com/raw")
	if err != nil {
		return ""
	}
	defer resp.Body.Close()
	content, _ := ioutil.ReadAll(resp.Body)
	return string(content)
}

func random(min, max int) int {
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, 4)
	res := (int(b[0]) << 24) + (int(b[1]) << 16) + (int(b[2]) << 8) + int(b[3])
	return res%(max-min) + min
}
