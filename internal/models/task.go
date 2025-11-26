package models

import (
	"time"

	"gorm.io/gorm"
)

type Task struct {
	ID          uint `gorm:"primarykey"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
	ProjectID   *uint          `gorm:"index"`
	ParentID    *uint          `gorm:"index"`
	Title       string         `gorm:"not null"`
	Description string
	Status      string `gorm:"default:'pending'"` // pending, in_progress, done
	Priority    int    `gorm:"default:0"`         // 0=low, 1=medium, 2=high
	Printed     bool   `gorm:"default:false"`
	PrintedAt   *time.Time
	CompletedAt *time.Time

	// Relationships
	Project  *Project
	Parent   *Task
	Subtasks []Task `gorm:"foreignkey:parentID"`
}
