package models

import (
	"example/config"
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// var DB_NAME string
func NewGorm() *gorm.DB {
	cfg := config.GetConfig()
	db, err := gorm.Open(sqlite.Open(cfg.DB_NAME), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database")
	}
	return db
}
