package ui

import (
	"context"
	"fmt"

	"github.com/FerrarioDev/thermaltodo/internal/models"
	taskrepository "github.com/FerrarioDev/thermaltodo/internal/repository/task"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Form struct {
	title       textinput.Model
	description textarea.Model
	parentID    *uint
	repo        taskrepository.TaskRepository
}

type TaskCreatedMsg struct {
	Task models.Task
}

type TaskCancelledMsg struct{}

func NewForm(repo taskrepository.TaskRepository, parentID *uint) *Form {
	form := &Form{
		repo:        repo,
		parentID:    parentID,
		title:       textinput.New(),
		description: textarea.New(),
	}
	form.title.Focus()
	return form
}

func (m Form) NewTask() tea.Msg {
	task := models.Task{
		Title:       m.title.Value(),
		Description: m.description.Value(),
		ParentID:    m.parentID,
		Status:      models.Todo,
	}
	_, err := m.repo.Create(context.Background(), &task)
	if err != nil {
		return fmt.Errorf("failed to create task: %w", err)
	}

	return TaskCreatedMsg{Task: task}
}

func (m Form) Init() tea.Cmd {
	return nil
}

func (m *Form) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, func() tea.Msg { return TaskCancelledMsg{} }
		case "tab":
			if m.title.Focused() {
				m.title.Blur()
				m.description.Focus()
				return m, textarea.Blink
			} else {
				m.description.Blur()
				m.title.Focus()
				return m, textinput.Blink
			}
		case "enter":
			if m.title.Focused() {
				m.title.Blur()
				m.description.Focus()
				return m, textarea.Blink
			}
			return m, func() tea.Msg { return m.NewTask() }
		}
	}
	// Update focused field
	var cmd tea.Cmd
	if m.title.Focused() {
		m.title, cmd = m.title.Update(msg)
		return m, cmd
	}
	m.description, cmd = m.description.Update(msg)
	return m, cmd
}

func (m Form) View() string {
	return lipgloss.JoinVertical(
		lipgloss.Left,
		"Create a new task",
		m.title.View(),
		m.description.View())
}
