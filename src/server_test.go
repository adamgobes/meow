package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func getRouter() *gin.Engine {
	r := gin.Default()
	return r
}

func getRegistrationPOSTPayload() string {
	params := url.Values{}
	params.Add("username", "u1")
	params.Add("password", "p1")

	return params.Encode()
}

func TestGetContactsEndpoiont(t *testing.T) {
	w := httptest.NewRecorder()

	r := getRouter()

	r.GET("/contacts", handleGetContacts)

	req, _ := http.NewRequest("GET", "/contacts", nil)

	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fail()
	}
}
