package handlers

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var Db *gorm.DB

func init() {

	var err error

	dsn := GetEnvVariable(DATABASE_URL)
	Db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {

		panic("failed to connect database")

	}

	// defer db.Close()
}
