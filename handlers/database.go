package handlers

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func init() {

	var err error
	dsn := "host=localhost port=5432 user=thorduragustsson dbname=all_quotes sslmode=disable password="
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	// db, err = gorm.Open("postgres", "host=localhost port=5432 user=thorduragustsson dbname=all_quotes sslmode=disable password=")
	if err != nil {

		panic("failed to connect database")

	}
	// db.LogMode(true)

	// defer db.Close()

	db.AutoMigrate(&Authors{})

	db.AutoMigrate(&Quotes{})
}
