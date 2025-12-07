package ui

import (
	"context"
	"fmt"

	"github.com/FerrarioDev/thermaltodo/internal/models"
	taskrepository "github.com/FerrarioDev/thermaltodo/internal/repository/task"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type Form struct {
	title       textinput.Model
	description textarea.Model
	parentID    *uint
	repo        taskrepository.TaskRepository
}

func NewForm(repo taskrepository.TaskRepository) *Form {
	form := &Form{repo: repo}
	form.title = textinput.New()
	form.title.Focus()
	form.description = textarea.New()
	return form
}

func (m Form) NewTask() tea.Msg {
	task := models.Task{
		Title:       m.title.Value(),
		Description: m.title.Value(),
		ParentID:    m.parentID,
		Status:      models.Todo,
	}
	_, err := m.repo.Create(context.Background(), &task)
	if err != nil {
		return fmt.Sprint(err)
	}

	return task
}

func (m Form) Init() tea.Cmd {
	return nil
}

func (m Form) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "enter":
			if m.title.Focused() {
				m.title.Blur()
				m.description.Focus()
				return m, textarea.Blink
			} else {
				pages[TaskForm] = m
				return pages[Board], m.NewTask
			}
		}
	}
	if m.title.Focused() {
		m.title, cmd = m.title.Update(msg)
		return m, cmd
	} else {
		m.description, cmd = m.description.Update(msg)
	}

	return m, cmd
}

func (m Form) View() string {
	return ""
}
