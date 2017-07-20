package model

import (
	"encoding/binary"
	"encoding/json"

	"github.com/boltdb/bolt"
)

const tmpFileBucketName = "tmpfile"

type TmpFile struct {
	ID       uint64 `json:"id"`
	Filename string `json:"filename"`
	Payload  []byte `json:"payload"`
}

func init() {
	openDB().Update(func(tx *bolt.Tx) error {
		var err error
		_, err = tx.CreateBucketIfNotExists([]byte(tmpFileBucketName))
		if err != nil {
			panic(err)
		}

		return nil
	})
}

func CreateTmpfile(filename string, payload []byte) (*TmpFile, error) {
	var tf *TmpFile
	err := openDB().Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(tmpFileBucketName))

		id, err := bucket.NextSequence()
		if err != nil {
			return err
		}

		idBytes := make([]byte, 8)
		binary.BigEndian.PutUint64(idBytes, id)

		tf = &TmpFile{
			ID:       id,
			Filename: filename,
			Payload:  payload,
		}

		b, err := json.Marshal(tf)
		if err != nil {
			return err
		}

		return bucket.Put(idBytes, b)
	})

	return tf, err
}

func GetTmpfileById(id uint64) (*TmpFile, error) {
	var tf TmpFile
	err := openDB().Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(tmpFileBucketName))

		idBytes := make([]byte, 8)
		binary.BigEndian.PutUint64(idBytes, id)

		b := bucket.Get(idBytes)
		return json.Unmarshal(b, &tf)
	})

	return &tf, err
}

func DelTmpFileById(id uint64) error {
	return openDB().Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(tmpFileBucketName))

		idBytes := make([]byte, 8)
		binary.BigEndian.PutUint64(idBytes, id)

		return bucket.Delete(idBytes)
	})
}
