package global

import (
	"encoding/hex"
	"fmt"
	"os"
	"sync"

	"github.com/YouDad/blockchain/log"
	"github.com/YouDad/blockchain/types"
	"github.com/boltdb/bolt"
)

type IDatabase interface {
	Clear()
	Get(key interface{}) (value []byte)
	Set(key interface{}, value []byte)
	Delete(key interface{})
	Foreach(func(k, v []byte) bool)
}

var databaseName string
var instanceBoltDB *bolt.DB
var onceBoltDB sync.Once

func RegisterDatabase(dbName string) {
	databaseName = dbName
}

func getDatabase() *bolt.DB {
	onceBoltDB.Do(func() {
		if !IsDatabaseExists() {
			log.Errln("No existing blockchain found, create one to continue.")
		}

		var err error
		instanceBoltDB, err = bolt.Open(databaseName, 0600, nil)
		log.Err(err)
	})
	return instanceBoltDB
}

func CreateDatabase() {
	if IsDatabaseExists() {
		log.Errln("Blockchain existed, Create failed.")
	}

	db, err := bolt.Open(databaseName, 0600, nil)
	log.Err(err)
	db.Close()
}

func IsDatabaseExists() bool {
	_, err := os.Stat(databaseName)
	return !os.IsNotExist(err)
}

func interfaceToString(key interface{}) string {
	var keyString string
	switch v := key.(type) {
	case types.HashValue:
		keyString = hex.EncodeToString(v)
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
	case types.HashValue:
		keyBytes = v
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
