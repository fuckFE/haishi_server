package server

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCreate(t *testing.T) {
	resp := httptest.NewRecorder()

	s := `{"user": "yejiayu", "password": "password"}`
	req := httptest.NewRequest("POST", "/api/users", strings.NewReader(s))
	req.Header.Set("Content-Type", "application/json")

	GetMainEngine().ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fail()
	}
}

func TestLogin(t *testing.T) {
	resp := httptest.NewRecorder()

	s := `{"user": "yejiayu", "password": "password"}`
	req := httptest.NewRequest("POST", "/api/users/login", strings.NewReader(s))
	req.Header.Set("Content-Type", "application/json")

	GetMainEngine().ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fail()
	}
}
