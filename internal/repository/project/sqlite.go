package projectrepository

import (
	"context"
	"errors"
	"time"

	"github.com/FerrarioDev/thermaltodo/internal/models"
	"gorm.io/gorm"
)

type SqliteProjectRepository struct {
	db *gorm.DB
}

func NewSqliteProjectRepository(db *gorm.DB) ProjectRepository {
	return &SqliteProjectRepository{db}
}

func (r *SqliteProjectRepository) Create(ctx context.Context, project *models.Project) (uint, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	result := r.db.WithContext(ctx).Create(project)
	if result.Error != nil {
		return 0, result.Error
	}
	return project.ID, nil
}

func (r *SqliteProjectRepository) GetByID(ctx context.Context, id uint) (*models.Project, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	var project models.Project
	result := r.db.WithContext(ctx).First(&project, "id = ?", id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &project, nil
}

func (r *SqliteProjectRepository) GetAll(ctx context.Context) ([]models.Project, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	var projects []models.Project
	result := r.db.WithContext(ctx).Find(&projects)
	if result.Error != nil {
		return nil, result.Error
	}
	return projects, nil
}

func (r *SqliteProjectRepository) Delete(ctx context.Context, id uint) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	result := r.db.WithContext(ctx).Delete(&models.Project{}, id)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("project not found")
	}

	return nil
}

func (r *SqliteProjectRepository) Update(ctx context.Context, project *models.Project) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if project.ID == 0 {
		return errors.New("invalid project id")
	}

	// Check if project exists first
	var existing models.Project
	if err := r.db.WithContext(ctx).First(&existing, project.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("project not found")
		}
		return err
	}

	// Update all fields except ID and timestamps (handled by GORM hooks)
	result := r.db.WithContext(ctx).Model(&models.Project{}).
		Where("id = ?", project.ID).
		Updates(project)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("no rows affected")
	}

	return nil
}
