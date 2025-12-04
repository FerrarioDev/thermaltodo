// Package database is for connection and database configuration
package database

import (
	"github.com/FerrarioDev/thermaltodo/internal/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func InitDB(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	if err := db.AutoMigrate(&models.Project{}, &models.Task{}); err != nil {
		return nil, err
	}

	return db, nil
}
