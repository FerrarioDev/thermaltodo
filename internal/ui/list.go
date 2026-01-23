package ui

import (
	"context"
	"fmt"

	"github.com/FerrarioDev/thermaltodo/internal/models"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const divisor = 4

// columnStyle = lipgloss.NewStyle().
//
//	Padding(1, 2)
var focusedStyle = lipgloss.NewStyle().
	Padding(1, 2).
	Border(lipgloss.RoundedBorder()).
	BorderForeground(lipgloss.Color("62"))

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

	tasks, _ = m.task.GetByParentID(context.Background(), m.currentParent, models.Todo)

	items := make([]list.Item, len(tasks))
	for i, task := range tasks {
		if len(task.Description) > 30 {
			task.Description = task.Description[:30] + "..."
		}
		items[i] = task.ConvertToUI()
	}

	m.list.SetItems(items)
}

func (m *App) navigateToChildren(task models.UITask) {
	_, err := m.task.GetByParentID(context.Background(), &task.ID, models.Todo)
	if err != nil {
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

func (m App) deleteTask(taskID uint) tea.Msg {
	err := m.task.Delete(context.Background(), taskID)
	if err != nil {
		return fmt.Errorf("failed to delete task: %w", err)
	}
	return TaskDeletedMsg{TaskID: taskID}
}

func (m App) printTask(selectedItem models.UITask) tea.Msg {
	ctx := context.Background()
	task, err := m.task.GetByID(ctx, selectedItem.ID)
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

	return TaskPrintedMsg{Task: *task}
}

func (m App) getAllDescendants(parentID *uint) []models.Task {
	tasks, err := m.task.GetByParentID(context.Background(), parentID, models.Todo)
	if err != nil {
		fmt.Printf("failed to get tasks: %v", err)
		return []models.Task{}
	}

	var allTasks []models.Task
	for _, task := range tasks {
		allTasks = append(allTasks, task)
		// Recursively get descendants of this task
		descendants := m.getAllDescendants(&task.ID)
		allTasks = append(allTasks, descendants...)
	}

	return allTasks
}

func (m App) printChildrens(parentID uint) tea.Msg {
	// Get all descendants (children, grandchildren, etc.)
	allTasks := m.getAllDescendants(&parentID)

	for _, task := range allTasks {
		job := models.PrintJob{
			TaskID:      task.ID,
			Title:       task.Title,
			Description: task.Description,
			CreatedAt:   task.CreatedAt,
		}
		if err := m.queue.Enqueue(job); err != nil {
			fmt.Printf("failed to enqueue job: %v", err)
		}
	}

	return TaskChildrenPrintedMsg{ParentID: parentID}
}

func (m App) completeTask(id uint) tea.Msg {
	if err := m.task.MarkComplete(context.Background(), id); err != nil {
		fmt.Printf("failed to MarkComplete")
	}

	return TaskCompletedMsg{id}
}
