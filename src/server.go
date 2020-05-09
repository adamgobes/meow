package main

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-contrib/cors"
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

type ContactWithPreview struct {
	UserID         string `json:"userId"`
	MessagePreview string `json:"messagePreview"`
}

func appendNoDup(contacts []string, target string) []string {
	for _, contact := range contacts {
		if contact == target {
			return contacts
		}
	}
	return append(contacts, target)
}

func handleGetContacts(c *gin.Context) {
	var messages []Message
	var contactIDs []string

	loggedInUser, err := unsignJWT(c)

	if err != nil {
		return
	}

	var allContacts []ContactWithPreview

	// find all messages that the logged in user participated in
	db.Where(&Message{Sender: loggedInUser}).Or(&Message{Receiver: loggedInUser}).Find(&messages)

	// get a unique list of all contacts
	for _, message := range messages {
		if loggedInUser == message.Sender {
			contactIDs = appendNoDup(contactIDs, message.Receiver)
		} else {
			contactIDs = appendNoDup(contactIDs, message.Sender)
		}
	}

	// for each one of those contacts, find the latest message exchanged between the logged in user and that contact
	for _, contactID := range contactIDs {
		db.Where(&Message{Sender: contactID, Receiver: loggedInUser}).Or(&Message{Sender: loggedInUser, Receiver: contactID}).Order("created_at desc").Order(1).Find(&messages)
		contactWithPreview := ContactWithPreview{UserID: contactID, MessagePreview: messages[0].Content}
		allContacts = append(allContacts, contactWithPreview)
	}

	c.JSON(http.StatusOK, gin.H{
		"data": allContacts,
	})
}

func setupRouter() *gin.Engine {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowMethods: []string{"GET", "POST", "HEAD", "PUT", "PATCH"},
		AllowHeaders: []string{"Origin", "Authorization"},
		AllowOriginFunc: func(origin string) bool {
			return origin == "http://localhost:3000" || origin == "https://mystatcat.com"
		},
	}))

	r.POST("/messages/new", handleWriteMessage)
	r.GET("/messages/:userId", handleGetMessages)
	r.GET("/contacts", handleGetContacts)

	return r
}

func runServer() error {
	router := setupRouter()

	return router.Run()
}
