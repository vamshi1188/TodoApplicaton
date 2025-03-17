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
	setupHeaders()

	// Populate the table with task data.
	refreshTable()

	// Make the table selectable and set up keybindings.
	table.SetSelectable(true, false).SetInputCapture(handleKeypress)

	// Create a header text view for the agenda title.
	title := tview.NewTextView().
		SetText("📝 My Terminal To-Do Agenda").
		SetTextAlign(tview.AlignCenter).
		SetTextColor(tcell.ColorGreen)

	// Create a footer text view for the help information.
	help := tview.NewTextView().
		SetText("A: Add | E: Edit | D: Delete | C: Complete | X: Pending | S: Save | H: Help | Q: Quit").
		SetTextAlign(tview.AlignCenter).
		SetTextColor(tcell.ColorYellow)

	// Arrange the title, table, and help in a vertical layout.
	layout := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(title, 1, 1, false). // Title at the top.
		AddItem(table, 0, 10, true). // Table in the middle.
		AddItem(help, 1, 1, false)   // Help at the bottom.

	// Run the application.
	if err := app.SetRoot(layout, true).Run(); err != nil {
		log.Fatalf("Error running application: %v", err)
	}
}

// Set up the table headers.
func setupHeaders() {
	headers := []string{"ID", "Task", "Due", "Priority", "Status"}
	for col, header := range headers {
		table.SetCell(0, col,
			tview.NewTableCell(header).
				SetTextColor(tcell.ColorYellow).
				SetAlign(tview.AlignCenter).
				SetSelectable(false))
	}
}

// Refresh the table with updated task data.
func refreshTable() {
	table.Clear()
	setupHeaders()

	for i, task := range tasks {
		row := i + 1
		table.SetCell(row, 0, tview.NewTableCell(fmt.Sprintf("%d", task.ID)).SetAlign(tview.AlignCenter))
		table.SetCell(row, 1, tview.NewTableCell(task.Description).SetAlign(tview.AlignLeft))
		table.SetCell(row, 2, tview.NewTableCell(task.Due.Format("2006-01-02")).SetAlign(tview.AlignCenter))
		table.SetCell(row, 3, tview.NewTableCell(task.Priority).SetAlign(tview.AlignCenter))
		table.SetCell(row, 4, tview.NewTableCell(task.Status).SetAlign(tview.AlignCenter))
	}
}

// Handle keypress events.
func handleKeypress(event *tcell.EventKey) *tcell.EventKey {
	switch event.Rune() {
	case 'c', 'C':
		toggleTaskStatus("✅ Done")
	case 'p', 'P':
		toggleTaskStatus("❌ Pending")
	case 'a', 'A':
		addTaskForm()
	case 'd', 'D':
		confirmDeleteTask()
	case 'e', 'E':
		editTaskForm()
	case 's', 'S':
		saveTasks()
		showMessage("Tasks saved successfully!")
	case 'h', 'H':
		showHelp()
	case 'q', 'Q':
		app.Stop()
	}
	return event
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
		AddDropDown("Priority", []string{"🔥 High", "👍 Medium", "⭐ Low"}, 0, nil)

	form.AddButton("Save", func() {
		desc := form.GetFormItem(0).(*tview.InputField).GetText()
		dueDateStr := form.GetFormItem(1).(*tview.InputField).GetText()
		dueDate, err := time.Parse("2006-01-02", dueDateStr)
		if err != nil || desc == "" {
			showError("Invalid input! Description cannot be empty & date must be YYYY-MM-DD.")
			return
		}
		_, priority := form.GetFormItem(2).(*tview.DropDown).GetCurrentOption()
		tasks = append(tasks, Task{taskIDCounter, desc, dueDate, priority, "❌ Pending"})
		taskIDCounter++
		refreshTable()
		app.SetRoot(table, true).SetFocus(table)
	})

	form.AddButton("Cancel", func() {
		app.SetRoot(table, true).SetFocus(table)
	})

	form.SetBorder(true).SetTitle("Add New Task").SetTitleAlign(tview.AlignLeft)
	app.SetRoot(form, true).SetFocus(form)
}

// Confirm deletion of a task.
func confirmDeleteTask() {
	row, _ := table.GetSelection()
	if row > 0 && row <= len(tasks) {
		modal := tview.NewModal().
			SetText("Are you sure you want to delete this task?").
			AddButtons([]string{"Yes", "No"}).
			SetDoneFunc(func(buttonIndex int, buttonLabel string) {
				if buttonLabel == "Yes" {
					tasks = append(tasks[:row-1], tasks[row:]...)
					refreshTable()
				}
				app.SetRoot(table, true).SetFocus(table)
			})
		app.SetRoot(modal, true)
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
			AddDropDown("Priority", []string{"🔥 High", "👍 Medium", "⭐ Low"}, 0, nil)

		form.AddButton("Save", func() {
			task.Description = form.GetFormItem(0).(*tview.InputField).GetText()
			dueDateStr := form.GetFormItem(1).(*tview.InputField).GetText()
			dueDate, err := time.Parse("2006-01-02", dueDateStr)
			if err != nil || task.Description == "" {
				showError("Invalid input! Description cannot be empty & date must be YYYY-MM-DD.")
				return
			}
			_, task.Priority = form.GetFormItem(2).(*tview.DropDown).GetCurrentOption()
			task.Due = dueDate
			refreshTable()
			app.SetRoot(table, true).SetFocus(table)
		})

		form.AddButton("Cancel", func() {
			app.SetRoot(table, true).SetFocus(table)
		})

		form.SetBorder(true).SetTitle("Edit Task").SetTitleAlign(tview.AlignLeft)
		app.SetRoot(form, true).SetFocus(form)
	}
}

// Save tasks to a JSON file.
func saveTasks() {
	file, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		log.Printf("Error marshalling tasks: %v", err)
		return
	}
	if err := os.WriteFile(dataFile, file, 0644); err != nil {
		log.Printf("Error saving tasks: %v", err)
	}
}

// Load tasks from a JSON file.
func loadTasks() {
	file, err := os.ReadFile(dataFile)
	if err != nil {
		return
	}
	if err := json.Unmarshal(file, &tasks); err != nil {
		log.Printf("Error loading tasks: %v", err)
	}
	if len(tasks) > 0 {
		taskIDCounter = tasks[len(tasks)-1].ID + 1
	}
}

// Show an error message in a modal.
func showError(message string) {
	modal := tview.NewModal().SetText(message).AddButtons([]string{"OK"}).SetDoneFunc(func(int, string) {
		app.SetRoot(table, true).SetFocus(table)
	})
	app.SetRoot(modal, true)
}

// Show a success message in a modal.
func showMessage(message string) {
	modal := tview.NewModal().SetText(message).AddButtons([]string{"OK"}).SetDoneFunc(func(int, string) {
		app.SetRoot(table, true).SetFocus(table)
	})
	app.SetRoot(modal, true)
}

// Show a help dialog with keybindings.
func showHelp() {
	modal := tview.NewModal().SetText("A: Add | E: Edit | D: Delete | C: Complete | X: Pending | S: Save | H: Help | Q: Quit").AddButtons([]string{"OK"}).SetDoneFunc(func(int, string) {
		app.SetRoot(table, true).SetFocus(table)
	})
	app.SetRoot(modal, true)
}