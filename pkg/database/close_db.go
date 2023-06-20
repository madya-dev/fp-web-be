package database

import (
	"gorm.io/gorm"
	"log"
)

func Close(db *gorm.DB) {
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("ERROR Close fatal error: %v", err)
	}

	err = sqlDB.Close()
	if err != nil {
		log.Fatalf("ERROR Close fatal error: %v", err)
	}
}
