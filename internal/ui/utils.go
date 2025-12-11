package ui

import "github.com/FerrarioDev/thermaltodo/internal/models"

type TaskCreatedMsg struct {
	Task models.Task
}

type TaskCancelledMsg struct{}
