// Package database is for connection and database configuration
package database

import (
	"fmt"

	"github.com/FerrarioDev/thermaltodo/internal/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func InitDB(dsn string) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Printf("failed to connect to database: %v", err)
		return nil
	}
	if err := db.AutoMigrate(&models.Project{}, &models.Task{}); err != nil {
		fmt.Printf("failed to automigrate: %v", err)
	}

	return db
}
