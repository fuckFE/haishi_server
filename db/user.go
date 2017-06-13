package db

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
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

const userBucketName = "users"
const (
	Root = iota
	Admin
)

func CreateUser(name, password string) (*User, error) {
	db := openDB()

	var u *User
	err := db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(userBucketName))
		if err != nil {
			return fmt.Errorf("CreateUser: %s", err)
		}

		u = &User{
			Name:     name,
			Password: hashPwd(password),
			Role:     Admin,
		}
		buf, err := json.Marshal(u)
		if err != nil {
			return fmt.Errorf("CreateUser: %s", err)
		}

		return b.Put([]byte(name), buf)
	})

	if err != nil {
		return nil, err
	}

	return u, nil
}

func Login(name, password string) bool {
	db := openDB()

	var u User
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(userBucketName))

		buf := b.Get([]byte(name))

		if err := json.Unmarshal(buf, &u); err != nil {
			return fmt.Errorf("Login: %s", err)
		}
		return nil
	})

	if err != nil {
		log.Printf("Login: %s\n", err)
		return false
	}

	if u.Password == hashPwd(password) {
		return true
	}
	return false
}

func hashPwd(password string) string {
	hashPwd := md5.Sum([]byte(password))
	return hex.EncodeToString(hashPwd[:])
}
