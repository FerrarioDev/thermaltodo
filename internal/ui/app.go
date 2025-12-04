package ui

import (
	projectrepository "github.com/FerrarioDev/thermaltodo/internal/repository/project"
	taskrepository "github.com/FerrarioDev/thermaltodo/internal/repository/task"
)

type App struct {
	task    *taskrepository.SqliteTaskRepository
	project *projectrepository.ProjectRepository
}
