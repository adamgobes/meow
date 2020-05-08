package main

import (
	"bytes"
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

var testSender string = "sender_username"
var testReceiver string = "receiver_username"

var testContent string = "hello world"
var newMessageID uint

var tokenRef = jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
	"userId": testSender,
})

var senderToken, err = tokenRef.SignedString([]byte(os.Getenv("APP_SECRET")))

var r *gin.Engine

func init() {
	db, dbError = gorm.Open("postgres", psqlInfo)
	migrate()

	if dbError != nil {
		panic("failed to connect database")
	}

	r = setupRouter()
}

func TestNewMessageEndpoint(t *testing.T) {
	w := httptest.NewRecorder()

	values := map[string]string{"receiver": testReceiver, "content": testContent}

	jsonValue, _ := json.Marshal(values)

	req, _ := http.NewRequest("POST", "/messages/new", bytes.NewBuffer(jsonValue))

	req.Header.Add("Authorization", "Bearer "+senderToken)

	r.ServeHTTP(w, req)

	res := struct{ Data Message }{}

	json.NewDecoder(w.Body).Decode(&res)
	newMessageID = res.Data.ID

	if w.Code != http.StatusOK {
		t.Fail()
	}

	if res.Data.Receiver != testReceiver {
		t.Fail()
	}

	if res.Data.Content != testContent {
		t.Fail()
	}

}

func TestGetContactsEndpoint(t *testing.T) {
	w := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", "/contacts", nil)

	req.Header.Add("Authorization", "Bearer "+senderToken)

	r.ServeHTTP(w, req)

	res := struct{ Data []ContactWithPreview }{}

	json.NewDecoder(w.Body).Decode(&res)

	if w.Code != http.StatusOK {
		t.Fail()
	}

	if res.Data[0].UserID != testReceiver {
		t.Fail()
	}

	if res.Data[0].MessagePreview != testContent {
		t.Fail()
	}

	db.Where("id = ?", newMessageID).Unscoped().Delete(Message{})
	db.Where("user_id = ?", testSender).Unscoped().Delete(User{})
	db.Where("user_id = ?", testReceiver).Unscoped().Delete(User{})
}
