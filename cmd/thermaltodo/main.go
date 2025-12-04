package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/FerrarioDev/thermaltodo/internal/database"
	"github.com/FerrarioDev/thermaltodo/internal/models"
	"github.com/FerrarioDev/thermaltodo/internal/printer"
	taskrepository "github.com/FerrarioDev/thermaltodo/internal/repository/task"
)

func main() {
	db, err := database.InitDB("thermaltodo.db")
	if err != nil {
		log.Fatal(err)
	}
	repository := taskrepository.NewSqliteTaskRepository(db)
	ctx := context.Background()
	tasks, err := repository.GetAll(ctx)
	if err != nil {
		log.Fatal(err)
	}

	p, bc := printer.Connect()

	printer := printer.NewEscPos(p, bc)

	job := models.PrintJob{
		TaskID:      1,
		Title:       "text",
		Description: "desc",
		Priority:    "1",
		CreatedAt:   time.Now(),
	}
	printer.Print(ctx, &job)

	for i, task := range tasks {
		fmt.Printf("task %d: %v", i, task)
	}
}
