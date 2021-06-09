package handlers

import (
	"github.com/Skjaldbaka17/quotes-api/structs"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var Db *gorm.DB

func init() {

	var err error
	dsn := "host=localhost port=5432 user=thorduragustsson dbname=all_quotes sslmode=disable password="
	Db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	// db, err = gorm.Open("postgres", "host=localhost port=5432 user=thorduragustsson dbname=all_quotes sslmode=disable password=")
	if err != nil {

		panic("failed to connect database")

	}
	// db.LogMode(true)

	// defer db.Close()

	Db.AutoMigrate(&structs.Authors{})

	Db.AutoMigrate(&structs.Quotes{})
}
