package main

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/joho/godotenv"
	"log"
)

var db *gorm.DB
var dbError error

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	db, dbError = gorm.Open("postgres", psqlInfo)

	if dbError != nil {
		panic("failed to connect database")
	}

	defer db.Close()

	log.Fatal(runServer())
}
