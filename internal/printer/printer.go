// Package printer is the interface for the printer and for the queues
package printer

import (
	"context"

	"github.com/FerrarioDev/thermaltodo/internal/models"
)

type Printer interface {
	Print(ctx context.Context, job *models.PrintJob) error
}

type Queue interface {
	Enqueue(job models.PrintJob) error
	Worker(ctx context.Context)
	Shutdown()
}
