package db

import (
	"fmt"
	"testing"

	"github.com/boltdb/bolt"
)

func TestPutBucket(t *testing.T) {
	db := openDB()
	err := db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("book"))
		if err != nil {
			return fmt.Errorf("create id bucket:%s", err)
		}

		return b.Put([]byte("name"), []byte("book_name"))
	})

	if err != nil {
		t.Fatal(err)
	}

	err = db.View(func(tx *bolt.Tx) error {
		name := tx.Bucket([]byte("book")).Get([]byte("name"))
		if string(name) != "book_name" {
			return fmt.Errorf("name == %s", string(name))
		}

		return nil
	})

	if err != nil {
		t.Fatal(err)
	}
}
