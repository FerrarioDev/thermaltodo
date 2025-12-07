package main

import (
	"fmt"
	"log"
	"os"

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

	pages := []tea.Model{ui.NewApp(repository), ui.NewForm(repository)}
	app := pages[ui.Board]

	p := tea.NewProgram(app)

	if _, err := p.Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
