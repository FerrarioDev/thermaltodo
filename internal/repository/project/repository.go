// Package projectRepository handles project database excecutions
package projectrepository

import (
	"context"

	"github.com/FerrarioDev/thermaltodo/internal/models"
)

type ProjectRepository interface {
	Create(ctx context.Context, project *models.Project) (uint, error)
	GetAll(ctx context.Context) ([]models.Project, error)
	GetByID(ctx context.Context, id uint) (*models.Project, error)
	Delete(ctx context.Context, id uint) error
	Update(ctx context.Context, project *models.Project) error
}
