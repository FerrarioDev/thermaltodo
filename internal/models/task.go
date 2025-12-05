package models

import (
	"time"

	"gorm.io/gorm"
)

type Status int

const (
	Todo Status = iota
	InProgress
	Done
)

type Task struct {
	ID          uint   `gorm:"primarykey"`
	ParentID    *uint  `gorm:"index"`
	Title       string `gorm:"not null"`
	Description string
	Status      Status `gorm:"default:'0'"` // todo, in_progress, done
	Printed     bool   `gorm:"default:false"`

	// Time related
	PrintedAt   *time.Time
	CompletedAt *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`

	// Relationships
	Project  *Project
	Parent   *Task
	Subtasks []Task `gorm:"foreignkey:parentID"`
}
