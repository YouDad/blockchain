package utils

import (
	"github.com/boltdb/bolt"
	"log"
)

var (
	dbFile       = "blockchain.db"
	blocksBucket = []byte("blocks")
	lastest      = []byte("lastest")
)

type Database struct {
	*bolt.DB
}

func OpenDatabase() (db *Database, exists bool) {
	_db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		log.Panic(err)
	}

	db = &Database{_db}
	db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(blocksBucket)
		if bucket == nil {
			log.Println("New a Database...")
			_, err := tx.CreateBucket([]byte(blocksBucket))
			if err != nil {
				log.Panic(err)
			}
		}

		exists = bucket != nil

		return nil
	})
	return db, exists
}

func (db *Database) Close() {
	db.DB.Close()
}

func (db *Database) GetLastest() []byte {
	return db.Get(lastest)
}

func (db *Database) Get(key []byte) (value []byte) {
	err := db.View(func(tx *bolt.Tx) error {
		value = tx.Bucket(blocksBucket).Get(key)
		return nil
	})
	if err != nil {
		log.Panic(err)
	}

	return value
}

func (db *Database) SetLastest(key []byte, value []byte) {
	db.Set(key, value)
	db.Set(lastest, value)
}

func (db *Database) Set(key []byte, value []byte) {
	err := db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(blocksBucket)

		err := bucket.Put(key, value)
		if err != nil {
			log.Panic(err)
		}

		return nil
	})
	if err != nil {
		log.Panic(err)
	}
}
