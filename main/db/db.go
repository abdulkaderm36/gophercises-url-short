package db

import (
	"github.com/boltdb/bolt"
)

type DB struct {
    DB *bolt.DB
}

const BUCKET = "PathsToUrls"

func InitDB() *DB {
    db, err := bolt.Open("my.db", 0600, nil)
    if err != nil {
        panic(err)
    }

    db.Update(func(tx *bolt.Tx) error {
        _, err := tx.CreateBucket([]byte(BUCKET)) 
        if err != nil {
           panic(err) 
        }
        return nil
    })

    return &DB{
        DB: db,
    }
}

func (db *DB) InitData(pathsToUrls map[string]string) {
    db.DB.Update(func(tx *bolt.Tx) error {
        b := tx.Bucket([]byte(BUCKET))
        for k, v := range pathsToUrls {
            if err := b.Put([]byte(k), []byte(v)); err != nil {
                panic(err)
            }
        }
        return nil
    })
}

func (db *DB) Close() {
    db.DB.Update(func(tx *bolt.Tx) error {
        tx.DeleteBucket([]byte(BUCKET))   
        return nil
    })
}
