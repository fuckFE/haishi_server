package model

import (
	"encoding/binary"
	"path/filepath"

	"github.com/boltdb/bolt"
)

var db *bolt.DB
var testDB *bolt.DB

func openDB() *bolt.DB {
	if db == nil {
		var err error
		path, err := filepath.Abs("data.db")
		if err != nil {
			panic(err)
		}
		db, err = bolt.Open(path, 0600, nil)
		if err != nil {
			panic(err)
		}
	}

	return db
}

func openTestDB() *bolt.DB {
	if testDB == nil {
		var err error
		testDB, err = bolt.Open("data-test.db", 0600, nil)
		if err != nil {
			panic(err)
		}
	}

	return testDB
}

func itob(v int) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}
