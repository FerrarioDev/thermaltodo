package repository

import (
	"context"
	"fmt"
	"testing"

	"github.com/FerrarioDev/thermaltodo/internal/database"
	"github.com/FerrarioDev/thermaltodo/internal/models"
	repository "github.com/FerrarioDev/thermaltodo/internal/repository/task"
)

func TestTaskRepository(t *testing.T) {
	db, _ := database.InitDB("test.db")
	repository := repository.NewSqliteTaskRepository(db)
	ctx := context.Background()
	t.Run("create task", func(t *testing.T) {
		task := models.Task{
			Title:       "test",
			Description: "this is a test task",
			Priority:    0,
		}
		id, err := repository.Create(ctx, &task)
		if err != nil {
			t.Error(err)
		}

		created, err := repository.GetByID(ctx, id)
		if err != nil {
			t.Error(err)
		}

		fmt.Printf("Created task: %v", created)
		if created.Title != task.Title {
			t.Errorf("got %s, want %s", created.Title, task.Title)
		}
	})
}
