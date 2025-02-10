package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// Task represents a to-do item.
type Task struct {
	ID          int       `json:"id"`
	Description string    `json:"description"`
	Due         time.Time `json:"due"`
	Priority    string    `json:"priority"`
	Status      string    `json:"status"`
}

var tasks []Task
var taskIDCounter int = 1
var app *tview.Application
var table *tview.Table

const dataFile = "tasks.json"

func main() {
	// Load tasks from file.
	loadTasks()

	// Create a new Tview application.
	app = tview.NewApplication()

	// Create a table with borders.
	table = tview.NewTable().SetBorders(true)

	// Set up the header row.
	headers := []string{"ID", "Task", "Due", "Priority", "Status"}
	for col, header := range headers {
		table.SetCell(0, col,
			tview.NewTableCell(header).
				SetTextColor(tcell.ColorYellow).
				SetAlign(tview.AlignCenter).
				SetSelectable(false))
	}

	// Populate the table with task data.
	refreshTable()

	// Make the table selectable.
	table.SetSelectable(true, false)

	// Add input capture for keybindings.
	table.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Rune() {
		case 'c', 'C': // Mark task as complete.
			toggleTaskStatus("✅ Done")
		case 'x', 'X': // Mark task as pending.
			toggleTaskStatus("❌ Pending")
		case 'a', 'A': // Add a new task.
			addTaskForm()
		case 'd', 'D': // Delete a task.
			deleteTask()
		case 'e', 'E': // Edit a task.
			editTaskForm()
		case 's', 'S': // Save tasks to file.
			saveTasks()
		case 'h', 'H': // Show help.
			showHelp()
		case 'q', 'Q': // Quit the application.
			app.Stop()
		}
		return event
	})

	// Create a header text view for the agenda title.
	title := tview.NewTextView().
		SetText("📝 My Terminal To-Do Agenda (Press 'H' for Help)").
		SetTextAlign(tview.AlignCenter).
		SetTextColor(tcell.ColorGreen)

	// Arrange the title and table in a vertical layout.
	layout := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(title, 1, 1, false).
		AddItem(table, 0, 10, true)

	// Run the application.
	if err := app.SetRoot(layout, true).Run(); err != nil {
		log.Fatalf("Error running application: %v", err)
	}
}

// Refresh the table with updated task data.
func refreshTable() {
	table.Clear()
	headers := []string{"ID", "Task", "Due", "Priority", "Status"}
	for col, header := range headers {
		table.SetCell(0, col,
			tview.NewTableCell(header).
				SetTextColor(tcell.ColorYellow).
				SetAlign(tview.AlignCenter).
				SetSelectable(false))
	}

	// Populate the table with task data.
	for i, task := range tasks {
		row := i + 1 // Row 0 is header.
		table.SetCell(row, 0, tview.NewTableCell(fmt.Sprintf("%d", task.ID)).
			SetAlign(tview.AlignCenter))
		table.SetCell(row, 1, tview.NewTableCell(task.Description).
			SetAlign(tview.AlignLeft))
		table.SetCell(row, 2, tview.NewTableCell(task.Due.Format("2006-01-02")).
			SetAlign(tview.AlignCenter))
		table.SetCell(row, 3, tview.NewTableCell(task.Priority).
			SetAlign(tview.AlignCenter))
		table.SetCell(row, 4, tview.NewTableCell(task.Status).
			SetAlign(tview.AlignCenter))
	}
}

// Toggle task status between "✅ Done" and "❌ Pending".
func toggleTaskStatus(status string) {
	row, _ := table.GetSelection()
	if row > 0 && row <= len(tasks) {
		tasks[row-1].Status = status
		refreshTable()
	}
}

// Add a new task using a form.
func addTaskForm() {
	form := tview.NewForm().
		AddInputField("Description", "", 30, nil, nil).
		AddInputField("Due Date (YYYY-MM-DD)", "", 10, nil, nil).
		AddDropDown("Priority", []string{"🔥 High", "👍 Medium", "⭐ Low"}, 0, nil).
		AddButton("Save", func() {
			description := form.GetFormItem(0).(*tview.InputField).GetText()
			dueDateStr := form.GetFormItem(1).(*tview.InputField).GetText()
			dueDate, _ := time.Parse("2006-01-02", dueDateStr)
			_, priority := form.GetFormItem(2).(*tview.DropDown).GetCurrentOption()

			task := Task{
				ID:          taskIDCounter,
				Description: description,
				Due:         dueDate,
				Priority:    priority,
				Status:      "❌ Pending",
			}
			tasks = append(tasks, task)
			taskIDCounter++
			refreshTable()
			app.SetRoot(table, true)
		}).
		AddButton("Cancel", func() {
			app.SetRoot(table, true)
		})

	form.SetBorder(true).SetTitle("Add New Task").SetTitleAlign(tview.AlignLeft)
	app.SetRoot(form, true)
}

// Delete the selected task.
func deleteTask() {
	row, _ := table.GetSelection()
	if row > 0 && row <= len(tasks) {
		tasks = append(tasks[:row-1], tasks[row:]...)
		refreshTable()
	}
}

// Edit the selected task using a form.
func editTaskForm() {
	row, _ := table.GetSelection()
	if row > 0 && row <= len(tasks) {
		task := tasks[row-1]
		form := tview.NewForm().
			AddInputField("Description", task.Description, 30, nil, nil).
			AddInputField("Due Date (YYYY-MM-DD)", task.Due.Format("2006-01-02"), 10, nil, nil).
			AddDropDown("Priority", []string{"🔥 High", "👍 Medium", "⭐ Low"}, 0, nil).
			AddButton("Save", func() {
				task.Description = form.GetFormItem(0).(*tview.InputField).GetText()
				dueDateStr := form.GetFormItem(1).(*tview.InputField).GetText()
				task.Due, _ = time.Parse("2006-01-02", dueDateStr)
				_, task.Priority = form.GetFormItem(2).(*tview.DropDown).GetCurrentOption()
				refreshTable()
				app.SetRoot(table, true)
			}).
			AddButton("Cancel", func() {
				app.SetRoot(table, true)
			})

		form.SetBorder(true).SetTitle("Edit Task").SetTitleAlign(tview.AlignLeft)
		app.SetRoot(form, true)
	}
}

// Save tasks to a JSON file.
func saveTasks() {
	file, _ := json.MarshalIndent(tasks, "", "  ")
	_ = os.WriteFile(dataFile, file, 0644)
}

// Load tasks from a JSON file.
func loadTasks() {
	file, err := os.ReadFile(dataFile)
	if err != nil {
		return
	}
	_ = json.Unmarshal(file, &tasks)
	if len(tasks) > 0 {
		taskIDCounter = tasks[len(tasks)-1].ID + 1
	}
}

// Show a help dialog with keybindings.
func showHelp() {
	helpText := `Keybindings:
- A: Add a new task
- E: Edit the selected task
- D: Delete the selected task
- C: Mark task as complete
- X: Mark task as pending
- S: Save tasks to file
- H: Show this help dialog
- Q: Quit the application`

	modal := tview.NewModal().
		SetText(helpText).
		AddButtons([]string{"OK"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			app.SetRoot(table, true)
		})

	app.SetRoot(modal, true)
}
