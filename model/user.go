package model

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/boltdb/bolt"
)

type User struct {
	Name     string    `json:"name"`
	Password string    `json:"password"`
	Role     UserRole  `json:"role"`
	LastTime time.Time `json:"lastTime"`
}

type UserRole int

const (
	Root = iota
	Admin
)

func init() {
	openDB().Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("users"))
		if err != nil {
			panic(err)
		}

		return nil
	})
}

func CreateUser(name, password string) (*User, error) {
	u := &User{
		Name:     name,
		Password: hashPwd(password),
		Role:     Admin,
	}
	buf, err := json.Marshal(u)
	if err != nil {
		return nil, fmt.Errorf("CreateUser: %s", err)
	}

	err = openDB().Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("users"))
		if err := bucket.Put([]byte(name), buf); err != nil {
			return fmt.Errorf("CreateUser: %s", err)
		}

		return nil
	})

	return u, err
}

func Login(name, password string) bool {
	success := true
	err := openDB().Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("users"))
		buf := bucket.Get([]byte(name))

		var u User
		if err := json.Unmarshal(buf, &u); err != nil {
			return err
		}

		if u.Password == hashPwd(password) {
			return nil
		}

		return errors.New("login fail")
	})

	if err != nil {
		success = false
		log.Println(err)
	}

	return success
}

func GetUser(name string) (*User, error) {
	var u User
	err := openDB().Update(func(tx *bolt.Tx) error {
		buf := tx.Bucket([]byte("users")).Get([]byte(name))
		if err := json.Unmarshal(buf, &u); err != nil {
			return err
		}

		return nil
	})

	return &u, err
}

func hashPwd(password string) string {
	hashPwd := md5.Sum([]byte(password))
	return hex.EncodeToString(hashPwd[:])
}
