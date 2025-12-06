package ui

import (
	"context"
	"fmt"

	"github.com/FerrarioDev/thermaltodo/internal/models"
	taskrepository "github.com/FerrarioDev/thermaltodo/internal/repository/task"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type App struct {
	focused  models.Status
	task     taskrepository.TaskRepository
	list     list.Model
	loaded   bool
	quitting bool
}

func NewApp(task taskrepository.TaskRepository) tea.Model {
	return &App{task: task}
}

func (m App) Init() tea.Cmd {
	return nil
}

func (m App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		if !m.loaded {
			m.initList(msg.Width, msg.Height)
			m.loaded = true
		}
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.quitting = true
			return m, tea.Quit
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
		todoView := m.list.View()

		return lipgloss.JoinHorizontal(lipgloss.Left,
			todoView)
	} else {
		return "loading..."
	}
}

func (m *App) initList(width, height int) {
	defaultList := list.New([]list.Item{}, list.NewDefaultDelegate(), width, height)
	m.list = defaultList
	tasks, err := m.task.GetAll(context.Background())
	if err != nil {
		fmt.Errorf("failed to render tasks: %v", err)
	}

	items := make([]list.Item, len(tasks))
	for i, task := range tasks {
		items[i] = task.ConvertToUI()
	}
	m.list.SetItems(items)
}
