// Package ui handles all the tui logic and styling
package ui

import (
	"fmt"

	"github.com/FerrarioDev/thermaltodo/internal/models"
	taskrepository "github.com/FerrarioDev/thermaltodo/internal/repository/task"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type status int

const (
	Board status = iota
	TaskForm
)

type App struct {
	task     taskrepository.TaskRepository
	list     list.Model
	form     *Form
	loaded   bool
	quitting bool

	currentParent *uint
	currentView   status
	breadcrumb    []breadcrumbItem
}

type breadcrumbItem struct {
	parentID *uint
	taskName string
}

func NewApp(task taskrepository.TaskRepository) tea.Model {
	return &App{task: task}
}

func (m App) Init() tea.Cmd {
	return nil
}

func (m *App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		if !m.loaded {
			focusedStyle.Width(msg.Width / divisor)
			focusedStyle.Height(msg.Height - divisor)
			m.initList(msg.Width, msg.Height)
			m.loaded = true
		}
	case TaskCreatedMsg:
		m.currentView = Board
		m.loadTasksForCurrentLevel()
		return m, nil
	case TaskDeletedMsg: // Task was deleted
		m.loadTasksForCurrentLevel() // Reload the list
		return m, nil
	case TaskCancelledMsg: // User cancelled form
		m.currentView = Board
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.quitting = true
			return m, tea.Quit
		case "enter", "l":
			selectedItem := m.list.SelectedItem()
			if selectedItem != nil {
				task := selectedItem.(models.UITask)
				m.navigateToChildren(task)
			}
		case "h", "backspace":
			m.navigateBack()
		case "n":
			selectedItem := m.list.SelectedItem()
			var parentID *uint
			if selectedItem != nil {
				task := selectedItem.(models.UITask)
				parentID = &task.ID
			}
			m.form = NewForm(m.task, parentID) // Pass parentID
			m.currentView = TaskForm
			return m, m.form.Init()

		case "d":

			selectedItem := m.list.SelectedItem()
			task := selectedItem.(models.UITask)
			if selectedItem != nil {
				return m, func() tea.Msg { return m.deleteTask(task.ID) }
			}
		}
	}
	if m.currentView == TaskForm {
		_, cmd := m.form.Update(msg)
		return m, cmd
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
