package printer

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"strings"

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
	createdAt := job.CreatedAt.Format("2006-01-02 15:04:05")
	receipt := `
================================
  TASK: [%d]
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
	if err := s.printer.Printf(receipt, job.TaskID, job.Title, job.Description, job.Project, createdAt); err != nil {
		return fmt.Errorf("failed to print task: %v", err)
	}

	if err := s.printer.FeedLines(len(strings.Split(receipt, "\n"))); err != nil {
		return fmt.Errorf("failed to feed lines: %v", err)
	}

	err := s.printer.Cut()
	if err != nil {
		return err
	}

	cmd := exec.Command("lp", "-o", "raw")
	cmd.Stdin = s.buf
	_, err = cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to execute command: %v", err)
	}
	return nil
}

func Connect() (*escpos.Printer, *BufferCloser) {
	buf := &BufferCloser{new(bytes.Buffer)}

	printer := escpos.NewPrinter(buf)
	return &printer, buf
}
