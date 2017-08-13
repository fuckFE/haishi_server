package model

import (
	"encoding/binary"
	"encoding/json"
	"time"

	"github.com/boltdb/bolt"
)

type Book struct {
	ID        uint64    `json:"id"`
	Filename  string    `json:"filename"`
	PayloadID uint64    `json:"payloadID"`
	Types     []uint64  `json:"types"`
	Last      time.Time `json:"last"`
	Payload   string    `json:"payload"`
	IsGarbate bool      `json:"isGarbate"`
}

type RawBook struct {
	ID      uint64 `json:"id"`
	Payload string `json:"payload"`
}

const bookBucketName = "books"
const rawBookBucketName = "rawBook"

func init() {
	openDB().Update(func(tx *bolt.Tx) error {
		var err error
		_, err = tx.CreateBucketIfNotExists([]byte(bookBucketName))
		if err != nil {
			panic(err)
		}
		_, err = tx.CreateBucketIfNotExists([]byte(rawBookBucketName))
		if err != nil {
			panic(err)
		}
		return nil
	})
}

func CreateBook(filename string, payload []byte, types []uint64) error {
	return openDB().Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bookBucketName))

		// done := false
		// bucket.ForEach(func(k, v []byte) error {
		// 	var b Book
		// 	if err := json.Unmarshal(v, &b); err != nil {
		// 		return err
		// 	}

		// 	if b.Filename == filename {
		// 		b.Payload = string(payload)
		// 		b.Types = types

		// 		nb, err := json.Marshal(b)
		// 		if err != nil {
		// 			return err
		// 		}

		// 		if err := bucket.Put(k, nb); err != nil {
		// 			return err
		// 		}

		// 		done = true
		// 	}
		// 	return nil
		// })

		// if done {
		// 	return nil
		// }

		id, err := bucket.NextSequence()
		if err != nil {
			return err
		}

		idBytes := make([]byte, 8)
		binary.BigEndian.PutUint64(idBytes, id)

		payloadID, err := createRawBook(tx.Bucket([]byte(rawBookBucketName)), string(payload))
		if err != nil {
			return err
		}

		book := &Book{
			ID:        id,
			Filename:  filename,
			Types:     types,
			PayloadID: payloadID,
			Last:      time.Now(),
		}

		b, err := json.Marshal(book)
		if err != nil {
			return err
		}

		return bucket.Put(idBytes, b)
	})
}

func GetBookByGrabate() ([]*Book, error) {
	bs := make([]*Book, 0)
	err := openDB().Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bookBucketName))

		c := bucket.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			var b Book
			if err := json.Unmarshal(v, &b); err != nil {
				return err
			}

			if b.IsGarbate {
				bs = append(bs, &b)
			}
		}

		return nil
	})

	return bs, err
}

func GetBook(id uint64, filterGarbate bool) ([]*Book, error) {
	bs := make([]*Book, 0)
	err := openDB().Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bookBucketName))

		c := bucket.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			var b Book
			if err := json.Unmarshal(v, &b); err != nil {
				return err
			}

			if filterGarbate && b.IsGarbate {
				continue
			}
			for _, t := range b.Types {
				if t == id {
					bs = append(bs, &b)
					break
				}
			}
		}

		return nil
	})

	return bs, err
}

func GetBookByID(id uint64) (*Book, error) {
	var book Book
	err := openDB().Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bookBucketName))

		idBytes := make([]byte, 8)
		binary.BigEndian.PutUint64(idBytes, id)

		b := bucket.Get(idBytes)
		if err := json.Unmarshal(b, &book); err != nil {
			return err
		}

		rb, err := getRawBook(tx.Bucket([]byte(rawBookBucketName)), book.PayloadID)
		if err != nil {
			return err
		}
		book.Payload = rb.Payload
		return nil
	})

	return &book, err
}

func DelBook(id uint64) error {
	return openDB().Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bookBucketName))

		idBytes := make([]byte, 8)
		binary.BigEndian.PutUint64(idBytes, id)

		return bucket.Delete(idBytes)
	})
}

func Garbate(id uint64) error {
	return openDB().Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bookBucketName))

		idBytes := make([]byte, 8)
		binary.BigEndian.PutUint64(idBytes, id)

		b := bucket.Get(idBytes)
		var book Book
		if err := json.Unmarshal(b, &book); err != nil {
			return err
		}

		book.IsGarbate = !book.IsGarbate

		b, err := json.Marshal(book)
		if err != nil {
			return err
		}

		return bucket.Put(idBytes, b)
	})
}

func UpdatePayload(id uint64, payload []byte, filename string) error {
	return openDB().Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bookBucketName))

		idBytes := make([]byte, 8)
		binary.BigEndian.PutUint64(idBytes, id)

		b := bucket.Get(idBytes)
		var book Book
		if err := json.Unmarshal(b, &book); err != nil {
			return err
		}

		book.Payload = string(payload)
		book.Last = time.Now()
		book.Filename = filename

		b, err := json.Marshal(book)
		if err != nil {
			return err
		}

		return bucket.Put(idBytes, b)
	})
}

func createRawBook(bucket *bolt.Bucket, payload string) (uint64, error) {
	id, err := bucket.NextSequence()
	if err != nil {
		return 0, err
	}

	idBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(idBytes, id)

	rb := &RawBook{
		ID:      id,
		Payload: payload,
	}

	b, err := json.Marshal(rb)
	if err != nil {
		return 0, err
	}
	if err := bucket.Put(idBytes, b); err != nil {
		return 0, err
	}
	return id, nil
}

func getRawBook(bucket *bolt.Bucket, id uint64) (*RawBook, error) {
	var rb RawBook
	idBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(idBytes, id)

	b := bucket.Get(idBytes)
	if err := json.Unmarshal(b, &rb); err != nil {
		return nil, err
	}

	return &rb, nil
}

func updateRawBook(bucket *bolt.Bucket, id uint64, payload string) error {
	return openDB().Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(rawBookBucketName))

		idBytes := make([]byte, 8)
		binary.BigEndian.PutUint64(idBytes, id)

		b := bucket.Get(idBytes)
		var rb RawBook
		if err := json.Unmarshal(b, &rb); err != nil {
			return err
		}

		rb.Payload = payload

		b, err := json.Marshal(rb)
		if err != nil {
			return err
		}

		return bucket.Put(idBytes, b)
	})
}

func updateBookTypes(id uint64, ts []uint64) error {
	return openDB().Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bookBucketName))

		idBytes := make([]byte, 8)
		binary.BigEndian.PutUint64(idBytes, id)

		b := bucket.Get(idBytes)
		var book Book
		if err := json.Unmarshal(b, &book); err != nil {
			return err
		}
		book.Types = ts
		b, err := json.Marshal(book)
		if err != nil {
			return err
		}

		return bucket.Put(idBytes, b)
	})
}
