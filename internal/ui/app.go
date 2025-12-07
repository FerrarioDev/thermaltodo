// Package ui handles all the tui logic and styling
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

const divisor = 4

var (
	columnStyle = lipgloss.NewStyle().
			Padding(1, 2)
	focusedStyle = lipgloss.NewStyle().
			Padding(1, 2).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("62"))
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
		case "enter":
			selectedItem := m.list.SelectedItem()
			if selectedItem != nil {
				task := selectedItem.(models.UITask)
				m.navigateToChildren(task)
			}
		case "esc", "backspace":
			m.navigateBack()
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

func (m *App) initList(width, height int) {
	defaultList := list.New([]list.Item{}, list.NewDefaultDelegate(), width, height)
	defaultList.Title = "Root Tasks"
	m.list = defaultList
	m.currentParent = nil
	m.breadcrumb = []breadcrumbItem{{parentID: nil, taskName: "Root"}}

	// Load root tasks (tasks with no parent)
	m.loadTasksForCurrentLevel()
}

func (m *App) loadTasksForCurrentLevel() {
	var tasks []models.Task
	var err error
	if m.currentParent == nil {
		// Load root tasks
		tasks, err = m.task.GetByParentID(context.Background(), nil)
	} else {
		tasks, err = m.task.GetByParentID(context.Background(), m.currentParent)
	}

	if err != nil {
		fmt.Errorf("failed to load tasks: %v", err)
	}

	items := make([]list.Item, len(tasks))
	for i, task := range tasks {
		items[i] = task.ConvertToUI()
	}

	m.list.SetItems(items)
}

func (m *App) navigateToChildren(task models.UITask) {
	children, err := m.task.GetByParentID(context.Background(), &task.ID)
	if err != nil || len(children) == 0 {
		return
	}

	m.breadcrumb = append(m.breadcrumb, breadcrumbItem{
		parentID: &task.ID,
		taskName: task.Title(),
	})

	m.currentParent = &task.ID
	m.list.Title = fmt.Sprintf(task.Title())

	m.loadTasksForCurrentLevel()
}

func (m *App) navigateBack() {
	if len(m.breadcrumb) <= 1 {
		return
	}

	m.breadcrumb = m.breadcrumb[:len(m.breadcrumb)-1]
	lastItem := m.breadcrumb[len(m.breadcrumb)-1]
	m.currentParent = lastItem.parentID

	if m.currentParent == nil {
		m.list.Title = "Root"
	} else {
		m.list.Title = lastItem.taskName
	}

	m.loadTasksForCurrentLevel()
}

func (m *App) renderBreadcrumb() string {
	if len(m.breadcrumb) == 0 {
		return ""
	}

	breadcrumbStr := "Navigation: "
	for i, item := range m.breadcrumb {
		if i > 0 {
			breadcrumbStr += " > "
		}
		breadcrumbStr += item.taskName
	}

	breadcrumbStr += "\n[Enter: Open subtasks | Esc: Go back | Q: Quit]"
	return breadcrumbStr
}
