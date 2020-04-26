package global

import (
	"sync"

	"github.com/YouDad/blockchain/log"
	"github.com/boltdb/bolt"
)

type boltDB struct {
	currentBucket string
	mutex         *sync.Mutex
}

var onlyMutexBoltDB sync.Mutex

func (db *boltDB) db(group int) *bolt.DB {
	return getDatabase(group)
}

func getBoltDB(bucket string) IDatabase {
	return &boltDB{bucket, &onlyMutexBoltDB}
}

func (db *boltDB) lock() {
	db.mutex.Lock()
}

func (db *boltDB) unlock() {
	db.mutex.Unlock()
}

func (db *boltDB) Clear(group int) {
	db.lock()
	defer db.unlock()
	log.SetCallerLevel(2)
	log.Debugf("Clear %s[%d]", db.currentBucket, group)
	log.SetCallerLevel(0)
	log.Err(db.db(group).Update(func(tx *bolt.Tx) error {
		err := tx.DeleteBucket([]byte(db.currentBucket))
		if err == bolt.ErrBucketNotFound {
			err = nil
		}

		if err != nil {
			return err
		}

		_, err = tx.CreateBucketIfNotExists([]byte(db.currentBucket))

		return err
	}))
}

func (db *boltDB) Get(group int, key interface{}) (value []byte) {
	db.lock()
	defer db.unlock()
	db.db(group).View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(db.currentBucket))
		if bucket == nil {
			value = nil
		} else {
			value = bucket.Get(interfaceToBytes(key))
		}
		log.SetCallerLevel(4)
		log.Debugf("Get %s[%d] %s %d", db.currentBucket, group, interfaceToString(key), len(value))
		log.SetCallerLevel(0)
		return nil
	})
	return value
}

func (db *boltDB) GetWithoutLog(group int, key interface{}) (value []byte) {
	db.lock()
	defer db.unlock()
	db.db(group).View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(db.currentBucket))
		if bucket == nil {
			value = nil
		} else {
			value = bucket.Get(interfaceToBytes(key))
		}
		return nil
	})
	return value
}

func (db *boltDB) Set(group int, key interface{}, value []byte) {
	db.lock()
	defer db.unlock()
	log.Err(db.db(group).Update(func(tx *bolt.Tx) error {
		log.SetCallerLevel(4)
		log.Debugf("Set %s[%d] %s %d", db.currentBucket, group, interfaceToString(key), len(value))
		log.SetCallerLevel(0)
		bucket := tx.Bucket([]byte(db.currentBucket))
		if bucket == nil {
			var err error
			bucket, err = tx.CreateBucketIfNotExists([]byte(db.currentBucket))
			log.Err(err)
		}
		return bucket.Put(interfaceToBytes(key), value)
	}))
}

func (db *boltDB) Delete(group int, key interface{}) {
	db.lock()
	defer db.unlock()
	log.Err(db.db(group).Update(func(tx *bolt.Tx) error {
		log.SetCallerLevel(4)
		log.Debugf("Delete %s[%d] %s", db.currentBucket, group, interfaceToString(key))
		log.SetCallerLevel(0)
		return tx.Bucket([]byte(db.currentBucket)).Delete(interfaceToBytes(key))
	}))
}

func (db *boltDB) Foreach(group int, fn func(k, v []byte) bool) {
	db.lock()
	defer db.unlock()
	log.SetCallerLevel(2)
	log.Debugf("Foreach %s[%d]", db.currentBucket, group)
	log.SetCallerLevel(0)
	db.db(group).View(func(tx *bolt.Tx) error {
		cursor := tx.Bucket([]byte(db.currentBucket)).Cursor()
		for k, v := cursor.First(); k != nil; k, v = cursor.Next() {
			if !fn(k, v) {
				break
			}
		}
		return nil
	})
}
