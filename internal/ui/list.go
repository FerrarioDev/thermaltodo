package ui

import (
	"context"
	"fmt"

	"github.com/FerrarioDev/thermaltodo/internal/models"
	"github.com/charmbracelet/bubbles/list"
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
		task.Description = task.Description[:30] + "..."
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
	m.list.Title = fmt.Sprint(task.Title())

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
