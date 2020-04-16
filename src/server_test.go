package main

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func getRegistrationPOSTPayload() string {
	params := url.Values{}
	params.Add("username", "u1")
	params.Add("password", "p1")

	return params.Encode()
}

var TOKEN string = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VySWQiOiJjanl4cGZzd3IwMGd6MDcwOXdrMXg5cnc0IiwiaWF0IjoxNTg2NDYwODMwfQ.WpoAEeONH-4xQK4xzrFEbuhXH82NJ7tJEOYrtRcjRqI"

func TestGetContactsEndpoint(t *testing.T) {
	w := httptest.NewRecorder()

	r := setupRouter()
	mockedDB, _, _ := sqlmock.New()

	db, _ = gorm.Open("postgres", mockedDB)

	req, _ := http.NewRequest("GET", "/contacts", nil)

	req.Header.Add("Authorization", "Bearer "+TOKEN)

	r.ServeHTTP(w, req)

	expected := `{"data":[]}`

	if w.Body.String() == expected {
		t.Fail()
	}

	if w.Code != http.StatusOK {
		t.Fail()
	}
}
