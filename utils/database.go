package utils

import (
	"log"
	"os"

	"github.com/boltdb/bolt"
)

var (
	lastest  = []byte("lastest")
	ref      = 0
	globalDB *bolt.DB
)

type Database struct {
	*bolt.DB

	current_bucket []byte
}

func IsDatabaseExists(dbFile string) bool {
	_, err := os.Stat(dbFile)
	return !os.IsNotExist(err)
}

func OpenDatabase(dbFile string) *Database {
	var err error

	if ref == 0 {
		globalDB, err = bolt.Open(dbFile, 0600, nil)
		if err != nil {
			log.Panic(err)
		}
		ref++
	} else {
		ref++
	}

	return (&Database{DB: globalDB}).Blocks()
}

func (db *Database) Close() {
	ref--
	if ref == 0 {
		db.DB.Close()
		globalDB = nil
	}
}

func (db *Database) getBucket(str string) *Database {
	db.current_bucket = []byte(str)
	return db
}

func (db *Database) Blocks() *Database {
	return db.getBucket("blocks")
}

func (db *Database) UTXOSet() *Database {
	return db.getBucket("utxo_set")
}

func (db *Database) Clear() error {
	return db.Update(func(tx *bolt.Tx) error {
		err := tx.DeleteBucket(db.current_bucket)
		if err == bolt.ErrBucketNotFound {
			err = nil
		}
		if err != nil {
			return err
		}

		_, err = tx.CreateBucket(db.current_bucket)
		return err
	})
}

func (db *Database) GetLastest() []byte {
	return db.Get(lastest)
}

func (db *Database) Get(key []byte) (value []byte) {
	err := db.View(func(tx *bolt.Tx) error {
		value = tx.Bucket(db.current_bucket).Get(key)
		return nil
	})
	if err != nil {
		log.Panic(err)
	}

	return value
}

func (db *Database) SetLastest(key []byte, value []byte) error {
	err := db.Set(key, value)
	if err != nil {
		return err
	}

	return db.Set(lastest, value)
}

func (db *Database) Set(key []byte, value []byte) error {
	return db.Update(func(tx *bolt.Tx) error {
		return tx.Bucket(db.current_bucket).Put(key, value)
	})
}

func (db *Database) Delete(key []byte) error {
	return db.Update(func(tx *bolt.Tx) error {
		return tx.Bucket(db.current_bucket).Delete(key)
	})
}

func (db *Database) Foreach(fn func(k, v []byte) (isContinue bool)) error {
	return db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(db.current_bucket)
		cursor := b.Cursor()

		for k, v := cursor.First(); k != nil; k, v = cursor.Next() {
			if !fn(k, v) {
				break
			}
		}

		return nil
	})
}
