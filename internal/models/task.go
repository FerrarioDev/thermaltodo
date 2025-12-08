package models

import (
	"fmt"
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
	Status      Status `gorm:"default:0"` // todo, in_progress, done
	Printed     bool   `gorm:"default:false"`

	// Time related
	PrintedAt   *time.Time
	CompletedAt *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`

	// Relationships
	Parent   *Task
	Subtasks []Task `gorm:"foreignkey:parentID"`
}

type UITask struct {
	ID          uint
	title       string
	description string
}

func (t Task) ConvertToUI() UITask {
	ui := UITask{
		ID:          t.ID,
		title:       t.Title,
		description: t.Description,
	}

	return ui
}

func (t UITask) FilterValue() string {
	return fmt.Sprint(t.ID)
}

// Title is typically used for display
func (t UITask) Title() string {
	return t.title
}

// Description is typically used for display
func (t UITask) Description() string {
	return t.description // or any other field
}
