package main

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"log"
)

var db *gorm.DB
var dbError error

func main() {

	db, dbError = gorm.Open("postgres", psqlInfo)

	if dbError != nil {
		panic("failed to connect database")
	}

	defer db.Close()

	log.Fatal(runServer())
}
