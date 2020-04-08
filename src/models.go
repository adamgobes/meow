package main

import (
	"time"
)

type GormModel struct {
	ID        uint       `gorm:"primary_key" json:"id"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	DeletedAt *time.Time `json:"deletedAt"`
}

// Message : Struct representing message relation in DB
type Message struct {
	GormModel
	Sender   string `sql:"type:VARCHAR(255) REFERENCES users(user_id)" json:"sender"`
	Receiver string `sql:"type:VARCHAR(255) REFERENCES users(user_id)" json:"receiver"`
	Content  string `sql:"type:text; json:content" json:"content"`
}

type User struct {
	GormModel
	UserID string `gorm:"unique;not null" json:"user_id"`
}
