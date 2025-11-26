package printer

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"

	"github.com/FerrarioDev/thermaltodo/internal/models"
	"github.com/joeyak/go-escpos"
)

type BufferCloser struct {
	*bytes.Buffer
}

func (bc *BufferCloser) Close() error {
	return nil
}

type EscPos struct {
	printer *escpos.Printer
	buf     *BufferCloser
}

func NewEscPos(printer *escpos.Printer, buf *BufferCloser) Printer {
	return &EscPos{printer, buf}
}

func (s *EscPos) Print(ctx context.Context, job *models.PrintJob) error {
	var priority string
	switch job.Priority {
	case "0":
		priority = "LOW PRIORITY"
	case "1":
		priority = "MEDIUM PRIORITY"
	case "2":
		priority = "HIGH PRIORITY"
	}

	createdAt := job.CreatedAt.Format("2006-01-02 15:04:05")
	receipt := `
================================
  TASK: [%s]
================================

Title:
  %s

Description:
  %s

--------------------------------
Project: %s
Created: %s
Status: PENDING
--------------------------------

Crumple & toss when complete! 

================================
--------------------------------
    `
	s.printer.Printf(receipt, priority, job.Title, job.Description, job.Project, createdAt)
	s.printer.FeedLines(4)
	err := s.printer.Cut()
	if err != nil {
		return nil
	}
	cmd := exec.Command("lp", "-o", "raw")
	cmd.Stdin = s.buf
	_, err = cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to print %v", err)
	}
	return nil
}

func Connect() (*escpos.Printer, *BufferCloser) {
	buf := &BufferCloser{new(bytes.Buffer)}

	printer := escpos.NewPrinter(buf)
	return &printer, buf
}
