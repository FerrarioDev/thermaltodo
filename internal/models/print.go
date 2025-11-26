package models

import "time"

type PrintJob struct {
	TaskID      uint
	Title       string
	Description string
	Project     string
	Priority    string
	CreatedAt   time.Time
}
