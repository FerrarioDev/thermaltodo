package ui

import "github.com/FerrarioDev/thermaltodo/internal/models"

type TaskCreatedMsg struct {
	Task models.Task
}

type TaskPrintedMsg struct {
	Task models.Task
}

type TaskChildrenPrintedMsg struct {
	ParentID uint
}

type TaskCompletedMsg struct {
	id uint
}

type TaskCancelledMsg struct{}

type TaskDeletedMsg struct {
	TaskID uint
}
