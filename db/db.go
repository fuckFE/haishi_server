package db

import (
	"encoding/binary"

	"github.com/boltdb/bolt"
)

var db *bolt.DB
var testDB *bolt.DB

func openDB() *bolt.DB {
	if db == nil {
		var err error
		db, err = bolt.Open("data.db", 0600, nil)
		if err != nil {
			panic(err)
		}
	}

	return db
}

func openTestDB() *bolt.DB  {
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
