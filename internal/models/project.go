// Models defines the entities for the program
package models

import (
	"time"

	"gorm.io/gorm"
)

type Project struct {
	ID          uint `gorm:"primarykey"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
	Name        string         `gorm:"not null"`
	Description string
	Color       string
	Tasks       []Task `gorm:"foreignKey:ProjectID"`
}
