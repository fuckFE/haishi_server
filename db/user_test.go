package db

import (
	"testing"
)

func TestCreate(t *testing.T) {
	u, err := CreateUser("yejiayu", "test")
	if err != nil {
		t.Fatal(err)
	}

	if u.Name != "yejiayu" {
		t.Fail()
	}
}

func TestLogin(t *testing.T) {
	success := Login("yejiayu", "test")
	if !success {
		t.Fail()
	}
}
