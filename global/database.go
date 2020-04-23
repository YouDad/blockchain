package global

import (
	"encoding/hex"
	"fmt"
	"sync"

	"github.com/YouDad/blockchain/log"
	"github.com/YouDad/blockchain/types"
	"github.com/boltdb/bolt"
)

type IDatabase interface {
	Clear(group int)
	Get(group int, key interface{}) (value []byte)
	GetWithoutLog(group int, key interface{}) (value []byte)
	Set(group int, key interface{}, value []byte)
	Delete(group int, key interface{})
	Foreach(group int, fn func(k, v []byte) bool)
}

var (
	instanceBoltDB = make(map[int]*bolt.DB)
	onceBoltDB     = make(map[int]*sync.Once)
	mutexBoltDB    = sync.Mutex{}
)

func getDatabase(group int) *bolt.DB {
	mutexBoltDB.Lock()
	defer mutexBoltDB.Unlock()
	once, ok := onceBoltDB[group]
	if !ok {
		onceBoltDB[group] = &sync.Once{}
		once = onceBoltDB[group]
	}
	once.Do(func() {
		databaseName := fmt.Sprintf("blockchain%s-%d.db", Port, group)

		var err error
		instanceBoltDB[group], err = bolt.Open(databaseName, 0600, nil)
		log.Err(err)
	})
	return instanceBoltDB[group]
}

func interfaceToString(key interface{}) string {
	var keyString string
	switch v := key.(type) {
	case []byte:
		keyString = hex.EncodeToString(v)
	case types.HashValue:
		keyString = v.String()
	case string:
		keyString = v
	case int32:
		keyString = fmt.Sprint(v)
	case nil:
		return "nil"
	}
	return keyString
}

func interfaceToBytes(key interface{}) []byte {
	keyBytes := []byte{}
	switch v := key.(type) {
	case []byte:
		keyBytes = v
	case types.HashValue:
		keyBytes = []byte(v)
	case string:
		keyBytes = []byte(v)
	case int32:
		bytes := [4]byte{}
		for i := 0; i < 4; i++ {
			bytes[i] = byte(v >> (i * 8))
			if v < 256<<(i*8) {
				keyBytes = bytes[:i+1]
				break
			}
		}
	case nil:
		log.Warnln("key==nil")
		return nil
	}
	return keyBytes
}
