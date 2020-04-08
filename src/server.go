package main

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"strings"
)

type NewMessagePayload struct {
	Content  string
	Receiver string
}

func unsignJWT(c *gin.Context) (string, error) {
	authHeader := c.Request.Header.Get("Authorization")
	token := strings.Replace(authHeader, "Bearer ", "", 1)

	parsed, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("APP_SECRET")), nil
	})

	if err != nil || !parsed.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "unauthorized"})
		return "", fmt.Errorf("Error unsigning JWT")
	}

	claims, _ := parsed.Claims.(jwt.MapClaims)

	loggedInUser := fmt.Sprintf("%s", claims["userId"])

	return loggedInUser, nil
}

func handleGetMessages(c *gin.Context) {
	var messages []Message

	userID := c.Param("userId")

	loggedInUser, err := unsignJWT(c)

	if err != nil {
		return
	}

	db.Where(&Message{Sender: loggedInUser, Receiver: userID}).Find(&messages)

	c.JSON(http.StatusOK, gin.H{
		"data": messages,
	})
}

func handleWriteMessage(c *gin.Context) {
	var payload NewMessagePayload

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	loggedInUser, err := unsignJWT(c)

	if err != nil {
		return
	}

	sender := User{UserID: loggedInUser}
	receiver := User{UserID: payload.Receiver}

	db.Where(sender).First(&sender)
	db.Where(receiver).First(&receiver)

	if missing := db.NewRecord(sender); missing {
		db.Create(&sender)
	}

	if missing := db.NewRecord(receiver); missing {
		db.Create(&receiver)
	}

	message := Message{Sender: loggedInUser, Receiver: payload.Receiver, Content: payload.Content}

	db.Create(&message)

	c.JSON(http.StatusOK, gin.H{
		"data": message,
	})
}

func handleGetContacts(c *gin.Context) {
	var messages []Message

	loggedInUser, err := unsignJWT(c)

	if err != nil {
		return
	}

	db.Select("DISTINCT(receiver)").Where(&Message{Sender: loggedInUser}).Find(&messages)

	c.JSON(http.StatusOK, gin.H{
		"data": messages,
	})
}

func runServer() error {
	r := gin.Default()

	r.POST("/messages/new", handleWriteMessage)
	r.GET("/messages/:userId", handleGetMessages)
	r.GET("/contacts", handleGetContacts)

	return r.Run()
}
