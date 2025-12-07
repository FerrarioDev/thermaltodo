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

var pages []tea.Model

const (
	Board status = iota
	TaskForm
)

type App struct {
	focused       models.Status
	task          taskrepository.TaskRepository
	list          list.Model
	loaded        bool
	quitting      bool
	currentParent *uint
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
			pages[Board] = m
			return pages[TaskForm].Update(nil)
		}
	}
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m App) View() string {
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
