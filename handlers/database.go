package handlers

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func init() {

	var err error
	db, err = gorm.Open("postgres", "host=localhost port=5432 user=thorduragustsson dbname=all_quotes sslmode=disable password=")
	if err != nil {

		panic("failed to connect database")

	}
	db.LogMode(true)

	// defer db.Close()

	db.AutoMigrate(&Authors{})

	db.AutoMigrate(&Quotes{})
}
