package taskrepository

import (
	"context"
	"errors"
	"time"

	"github.com/FerrarioDev/thermaltodo/internal/models"
	"gorm.io/gorm"
)

const timeout time.Duration = 5

type SqliteTaskRepository struct {
	db *gorm.DB
}

func NewSqliteTaskRepository(db *gorm.DB) TaskRepository {
	return &SqliteTaskRepository{db}
}

func (r *SqliteTaskRepository) Create(ctx context.Context, task *models.Task) (uint, error) {
	ctx, cancel := context.WithTimeout(ctx, timeout*time.Second)
	defer cancel()

	if task.Title == "" {
		return 0, errors.New("missing task title")
	}

	// if err := r.db.WithContext(ctx).First(&models.Project{}, task.Project).Error; err != nil {
	//  if errors.Is(err, gorm.ErrRecordNotFound) {
	//     return 0, errors.New("project not found")
	//  }
	//  return 0, err
	//  }

	if err := r.db.WithContext(ctx).First(&models.Task{}, task.ParentID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, errors.New("parent not found")
		}
		return 0, err
	}

	result := r.db.WithContext(ctx).Create(task)
	if result.Error != nil {
		return 0, result.Error
	}
	return task.ID, nil
}

func (r *SqliteTaskRepository) GetByID(ctx context.Context, id uint) (*models.Task, error) {
	ctx, cancel := context.WithTimeout(ctx, timeout*time.Second)
	defer cancel()

	var task models.Task
	result := r.db.WithContext(ctx).First(&task, "id = ?", id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &task, nil
}

func (r *SqliteTaskRepository) GetAll(ctx context.Context) ([]models.Task, error) {
	ctx, cancel := context.WithTimeout(ctx, timeout*time.Second)
	defer cancel()
	var tasks []models.Task
	result := r.db.WithContext(ctx).Find(&tasks)
	if result.Error != nil {
		return nil, result.Error
	}
	return tasks, nil
}

// func (r *SqliteTaskRepository) GetByProject(ctx context.Context, projectID uint) ([]models.Task, error) {
// 	ctx, cancel := context.WithTimeout(ctx, timeout*time.Second)
// 	defer cancel()
// 	var tasks []models.Task
// 	result := r.db.WithContext(ctx).Where("project_id", projectID).Find(&tasks)
// 	if result.Error != nil {
// 		return nil, result.Error
// 	}
// 	return tasks, nil
// }

func (r *SqliteTaskRepository) GetPending(ctx context.Context) ([]models.Task, error) {
	ctx, cancel := context.WithTimeout(ctx, timeout*time.Second)
	defer cancel()
	var tasks []models.Task
	result := r.db.WithContext(ctx).Where("status = ?", "pending").Find(&tasks)
	if result.Error != nil {
		return nil, result.Error
	}
	return tasks, nil
}

func (r *SqliteTaskRepository) GetForPrinting(ctx context.Context) ([]models.Task, error) {
	ctx, cancel := context.WithTimeout(ctx, timeout*time.Second)
	defer cancel()
	var tasks []models.Task
	result := r.db.WithContext(ctx).Where("printed", false).Find(&tasks)
	if result.Error != nil {
		return nil, result.Error
	}
	return tasks, nil
}

func (r *SqliteTaskRepository) Delete(ctx context.Context, id uint) error {
	ctx, cancel := context.WithTimeout(ctx, timeout*time.Second)
	defer cancel()
	result := r.db.WithContext(ctx).Delete(&models.Task{}, id)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("task not found")
	}

	return nil
}

func (r *SqliteTaskRepository) Update(ctx context.Context, task *models.Task) error {
	ctx, cancel := context.WithTimeout(ctx, timeout*time.Second)
	defer cancel()

	if task.ID == 0 {
		return errors.New("invalid task id")
	}

	var existing models.Task
	if err := r.db.WithContext(ctx).First(&existing, task.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("task not found")
		}
		return err
	}

	result := r.db.WithContext(ctx).Model(&models.Task{}).
		Where("id = ?", task.ID).
		Updates(task)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("no rows affected")
	}

	return nil
}

func (r *SqliteTaskRepository) MarkComplete(ctx context.Context, id uint) error {
	ctx, cancel := context.WithTimeout(ctx, timeout*time.Second)
	defer cancel()

	var existing models.Task
	if err := r.db.WithContext(ctx).First(&existing, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("task not found")
		}
		return err
	}

	now := time.Now()
	existing.Status = models.Done
	existing.CompletedAt = &now

	if err := r.db.WithContext(ctx).Save(&existing).Error; err != nil {
		return err
	}
	return nil
}
