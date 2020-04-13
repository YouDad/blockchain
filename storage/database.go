package storage

import (
	"log"
	"os"

	"github.com/boltdb/bolt"
)

type Database interface {
	IsExists() bool
	SetTable(table string) Database
	Get(key string) (value string)
	Set(key string, value string)
}

var (
	databaseName string
	globalDB     *BoltDB
)

type BoltDB struct {
	*bolt.DB
	CurrentBucket []byte
}

func RegisterDatabase(dbName string) {
	databaseName = dbName
}

func GetDatabase() Database {
	return globalDB
}

func (db *BoltDB) IsExists() bool {
	if db != nil {
		return true
	} else {
		var err error
		db.DB, err = bolt.Open(databaseName, 0600, nil)
		if err != nil {
			log.Println(err)
			os.Exit(1)
			return false
		}
		return true
	}
}

func (db *BoltDB) SetTable(table string) Database {
	db.CurrentBucket = []byte(table)
	return db
}

func (db BoltDB) Get(key string) (value string) {
	err := db.View(func(tx *bolt.Tx) error {
		value = string(tx.Bucket(db.CurrentBucket).Get([]byte(key)))
		return nil
	})
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	return value
}

func (db BoltDB) Set(key string, value string) {
	err := db.View(func(tx *bolt.Tx) error {
		return tx.Bucket(db.CurrentBucket).Put([]byte(key), []byte(value))
	})
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
