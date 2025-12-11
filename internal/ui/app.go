// Package ui handles all the tui logic and styling
package ui

import (
	"context"
	"fmt"

	"github.com/FerrarioDev/thermaltodo/internal/models"
	"github.com/FerrarioDev/thermaltodo/internal/printer"
	taskrepository "github.com/FerrarioDev/thermaltodo/internal/repository/task"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type status int

// Model handler enum
const (
	Board status = iota
	TaskForm
)

// App model definition (Represents the main model for the UI)
type App struct {
	task     taskrepository.TaskRepository
	queue    printer.Queue
	list     list.Model
	form     *Form
	loaded   bool
	quitting bool

	// List Handlers
	currentParent *uint
	currentView   status
	breadcrumb    []breadcrumbItem
}

type breadcrumbItem struct {
	parentID *uint
	taskName string
}

/* Initialization */

func NewApp(task taskrepository.TaskRepository, queue printer.Queue) tea.Model {
	return &App{task: task, queue: queue}
}

func (m App) Init() tea.Cmd {
	return nil
}

func (m *App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	// Load the screen with the asign width and height of the terminal
	case tea.WindowSizeMsg:
		if !m.loaded {
			focusedStyle.Width(msg.Width / divisor)
			focusedStyle.Height(msg.Height - divisor)
			m.initList(msg.Width, msg.Height)
			m.loaded = true
		}

	/* Reload updated lists when tasks jobs where made */

	case TaskCreatedMsg: // User created task
		m.currentView = Board
		m.loadTasksForCurrentLevel()
		return m, nil
	case TaskDeletedMsg: // Task was deleted
		m.loadTasksForCurrentLevel()
		return m, nil
	case TaskCancelledMsg: // User cancelled task
		m.currentView = Board
		return m, nil

	// quit handler
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.quitting = true
			return m, tea.Quit
		}
	}

	if m.currentView == TaskForm {
		_, cmd := m.form.Update(msg)
		return m, cmd
	}

	// Board-specific key handlers
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		switch keyMsg.String() {

		// Move to next list
		case "enter", "l":
			selectedItem := m.list.SelectedItem()
			if selectedItem != nil {
				task := selectedItem.(models.UITask)
				m.navigateToChildren(task)
			}

		// Move to previous list
		case "h", "backspace":
			m.navigateBack()

		// Create new task
		case "n":
			m.form = NewForm(m.task, m.currentParent)
			m.currentView = TaskForm
			return m, m.form.Init()

		// Delete task
		case "d":
			selectedItem := m.list.SelectedItem()
			if selectedItem != nil {
				task := selectedItem.(models.UITask)
				return m, func() tea.Msg { return m.deleteTask(task.ID) }
			}
		case "p":
			selectedItem := m.list.SelectedItem().(models.UITask)

			task, err := m.task.GetByID(context.Background(), selectedItem.ID)
			if err != nil {
				fmt.Printf("failed to get task: %v", err)
			}

			job := models.PrintJob{
				TaskID:      task.ID,
				Title:       task.Title,
				Description: task.Description,
				CreatedAt:   task.CreatedAt,
			}

			if err := m.queue.Enqueue(job); err != nil {
				fmt.Printf("failed to enqueue: %v", err)
			}
		}
	}
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m App) View() string {
	if m.currentView == TaskForm {
		return m.form.View()
	}
	if m.quitting {
		return ""
	}
	if m.loaded {
		breadcrumbView := m.renderBreadcrumb()

		return fmt.Sprintf("%s\n\n%s", breadcrumbView, m.list.View())
	} else {
		return "loading..."
	}
}
