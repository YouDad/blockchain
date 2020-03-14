package store

import (
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"sync"

	"github.com/YouDad/blockchain/log"
	"github.com/YouDad/blockchain/types"
	"github.com/boltdb/bolt"
)

type IDatabase interface {
	SetTable(table string) IDatabase
	Clear()
	Get(key interface{}) (value []byte)
	Set(key interface{}, value []byte)
	Delete(key interface{})
	Foreach(func(k, v []byte) bool)
}

var (
	databaseName string
	opened       = false
	globalDB     = &BoltDB{}
	dbMutex      = sync.Mutex{}
)

type BoltDB struct {
	*bolt.DB
	CurrentBucket []byte
}

func RegisterDatabase(dbName string) {
	databaseName = dbName
}

func GetDatabase() IDatabase {
	if !IsDatabaseExists() {
		log.Errln(errors.New("No existing blockchain found, create one to continue."))
		return nil
	}

	if !opened {
		var err error
		globalDB.DB, err = bolt.Open(databaseName, 0600, nil)
		log.Err(err)
		opened = true
	}
	return globalDB
}

func CreateDatabase() IDatabase {
	if IsDatabaseExists() {
		log.Errln(errors.New("Blockchain existed, Create failed."))
		return nil
	}

	var err error
	globalDB.DB, err = bolt.Open(databaseName, 0600, nil)
	log.Err(err)
	opened = true
	return globalDB
}

func lock() {
	log.SetCallerLevel(2)
	log.Debugln("DB Lock")
	log.SetCallerLevel(0)
	dbMutex.Lock()
	// log.PrintStack()
}
func unlock() {
	dbMutex.Unlock()
	log.Debugln("DB Unlock")
}

func IsDatabaseExists() bool {
	_, err := os.Stat(databaseName)
	return !os.IsNotExist(err)
}

func (db *BoltDB) SetTable(table string) IDatabase {
	lock()
	db.CurrentBucket = []byte(table)
	return db
}

func (db *BoltDB) Clear() {
	log.Err(db.Update(func(tx *bolt.Tx) error {
		log.Debugf("Clear %s\n", db.CurrentBucket)
		log.PrintStack()
		err := tx.DeleteBucket(db.CurrentBucket)
		if err == bolt.ErrBucketNotFound {
			err = nil
		}
		if err != nil {
			return err
		}
		_, err = tx.CreateBucketIfNotExists(db.CurrentBucket)
		return err
	}))
	unlock()
}

func InterfaceToString(key interface{}) string {
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

func InterfaceToBytes(key interface{}) []byte {
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

func (db BoltDB) Get(key interface{}) (value []byte) {
	log.Err(db.View(func(tx *bolt.Tx) error {
		value = tx.Bucket(db.CurrentBucket).Get(InterfaceToBytes(key))
		k, o := key.(string)
		if len(value) == 0 && o && k == "lastest" {
			log.PrintStack()
		}
		log.SetCallerLevel(3)
		log.Debugf("Get %s %s %s\n", string(db.CurrentBucket), InterfaceToString(key), value)
		log.SetCallerLevel(0)
		return nil
	}))
	unlock()
	return value
}

func (db BoltDB) Set(key interface{}, value []byte) {
	log.Err(db.Update(func(tx *bolt.Tx) error {
		log.SetCallerLevel(3)
		log.Debugf("Set %s %s %s\n", string(db.CurrentBucket), InterfaceToString(key), value)
		log.SetCallerLevel(0)
		return tx.Bucket(db.CurrentBucket).Put(InterfaceToBytes(key), value)
	}))
	unlock()
}

func (db *BoltDB) Delete(key interface{}) {
	log.Err(db.Update(func(tx *bolt.Tx) error {
		log.SetCallerLevel(3)
		log.Debugln("Delete", string(db.CurrentBucket), InterfaceToString(key))
		log.SetCallerLevel(0)
		return tx.Bucket(db.CurrentBucket).Delete(InterfaceToBytes(key))
	}))
	unlock()
}

func (db BoltDB) Foreach(fn func(k, v []byte) bool) {
	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(db.CurrentBucket)
		cursor := b.Cursor()

		for k, v := cursor.First(); k != nil; k, v = cursor.Next() {
			if !fn(k, v) {
				break
			}
		}
		return nil
	})
	unlock()
}
