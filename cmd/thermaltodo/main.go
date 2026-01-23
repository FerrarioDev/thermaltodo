package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/FerrarioDev/thermaltodo/internal/database"
	"github.com/FerrarioDev/thermaltodo/internal/printer"
	taskrepository "github.com/FerrarioDev/thermaltodo/internal/repository/task"
	"github.com/FerrarioDev/thermaltodo/internal/ui"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	dbPath := filepath.Join(os.Getenv("HOME"), ".thermaltodo", "thermaltodo.db")
	if err := os.MkdirAll(filepath.Dir(dbPath), 0o755); err != nil {
		fmt.Println("failed to create db")
	}
	db, err := database.InitDB(dbPath)
	if err != nil {
		log.Fatal(err)
	}
	repository := taskrepository.NewSqliteTaskRepository(db)
	e, bc := printer.Connect()

	escpos := printer.NewEscPos(e, bc)
	queue := printer.NewPrintQueue(escpos, 4)
	queue.Start(context.Background(), 4)

	app := ui.NewApp(repository, queue)

	p := tea.NewProgram(app)

	if _, err := p.Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
