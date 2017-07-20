package model

import (
	"encoding/binary"
	"encoding/json"

	"github.com/boltdb/bolt"
)

type Type struct {
	ID       uint64 `json:"id"`
	Name     string `json:"name"`
	Category uint   `json:"category"` // 1主题 2类型
}

const typeBucketName = "types"

func init() {
	openDB().Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(typeBucketName))
		if err != nil {
			panic(err)
		}

		return nil
	})
}

func CreateType(name string, category uint) (*Type, error) {
	var t *Type
	err := openDB().Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(typeBucketName))
		id, err := bucket.NextSequence()
		if err != nil {
			return err
		}

		idBytes := make([]byte, 8)
		binary.BigEndian.PutUint64(idBytes, id)
		t = &Type{
			ID:       id,
			Name:     name,
			Category: category,
		}
		b, err := json.Marshal(t)
		if err != nil {
			return err
		}
		return bucket.Put(idBytes, b)
	})

	return t, err
}

func GetTypes() ([]Type, error) {
	ts := make([]Type, 0)
	err := openDB().View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(typeBucketName))
		c := bucket.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			var t Type
			if err := json.Unmarshal(v, &t); err != nil {
				return err
			}
			ts = append(ts, t)
		}

		return nil
	})

	return ts, err
}
