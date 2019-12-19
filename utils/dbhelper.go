package utils

import (
	"log"
	"os"

	"github.com/boltdb/bolt"
)

var (
	dbFile       = "blockchain.db"
	blocksBucket = []byte("blocks")
	lastest      = []byte("lastest")
)

type Database struct {
	*bolt.DB
}

func IsDatabaseExists() bool {
	_, err := os.Stat(dbFile)
	return !os.IsNotExist(err)
}

func CreateDatabase() *Database {
	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		log.Panic(err)
	}

	db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucket(blocksBucket)
		return err
	})

	return &Database{db}
}

func OpenDatabase() (db *Database) {
	_db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		log.Panic(err)
	}

	return &Database{_db}
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
