package printer

import (
	"context"
	"fmt"
	"sync"

	"github.com/FerrarioDev/thermaltodo/internal/models"
)

type PrinterQueue struct {
	jobs     chan models.PrintJob
	printer  Printer
	wg       sync.WaitGroup
	stopChan chan struct{}
	mu       sync.Mutex
}

func NewPrintQueue(printer Printer, workers int) Queue {
	return &PrinterQueue{
		jobs:     make(chan models.PrintJob, workers),
		printer:  printer,
		wg:       sync.WaitGroup{},
		stopChan: make(chan struct{}),
		mu:       sync.Mutex{},
	}
}

func (pq *PrinterQueue) Enqueue(job models.PrintJob) error {
	select {
	case pq.jobs <- job:
		return nil
	case <-pq.stopChan:
		return fmt.Errorf("queue is shutting down")
	default:
		return fmt.Errorf("queue is full")
	}
}

func (pq *PrinterQueue) Worker(ctx context.Context) {
	// Worker goroutine that processes jobs
	defer pq.wg.Done() // Decrement when this worker exits

	for {
		select {
		case job, ok := <-pq.jobs: // Wait for job
			if !ok {
				// Channel closed, no more jobs
				return
			}

			// Process the job
			if err := pq.printer.Print(ctx, &job); err != nil {
				// Log error (in real code, use proper logging)
				fmt.Printf("Print error: %v\n", err)
			}

		case <-ctx.Done(): // Context cancelled
			return
		}
	}
}

func (pq *PrinterQueue) Shutdown() {
	// Graceful Shutdown: drain queue, wait for workers
	pq.mu.Lock()
	defer pq.mu.Unlock()

	close(pq.stopChan)

	close(pq.jobs)

	pq.wg.Wait()
}

func (pq *PrinterQueue) Start(ctx context.Context, workers int) {
	for range workers {
		pq.wg.Add(1)
		go pq.Worker(ctx)
	}
}
