package taskrepository

import (
	"context"

	"github.com/FerrarioDev/thermaltodo/internal/models"
)

type TaskRepository interface {
	Create(ctx context.Context, task *models.Task) (uint, error)
	GetByID(ctx context.Context, id uint) (*models.Task, error)
	GetAll(ctx context.Context) ([]models.Task, error)
	GetByProject(ctx context.Context, projectID uint) ([]models.Task, error)
	GetPending(ctx context.Context) ([]models.Task, error)
	GetForPrinting(ctx context.Context) ([]models.Task, error)
	Update(ctx context.Context, task *models.Task) error
	Delete(ctx context.Context, id uint) error
	MarkComplete(ctx context.Context, id uint) error
}
