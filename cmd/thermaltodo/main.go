package main

import (
	"log"

	"github.com/FerrarioDev/thermaltodo/internal/database"
	taskrepository "github.com/FerrarioDev/thermaltodo/internal/repository/task"
	"github.com/FerrarioDev/thermaltodo/internal/ui"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	db, err := database.InitDB("thermaltodo.db")
	if err != nil {
		log.Fatal(err)
	}
	repository := taskrepository.NewSqliteTaskRepository(db)

	app := ui.NewApp(repository)

	p := tea.NewProgram(app)

	p.Run()
}
