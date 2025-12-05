package ui

import (
	taskrepository "github.com/FerrarioDev/thermaltodo/internal/repository/task"
)

type App struct {
	task *taskrepository.SqliteTaskRepository
}
