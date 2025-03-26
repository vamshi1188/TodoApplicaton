package main

import (
		"time"
		"github.com/rivo/tview"
		"fmt"
)

type Task struct {
	ID          int       `json:"id"`
	Description string    `json:"description"`
	Due         time.Time `json:"due"`
	Priority    string    `json:"priority"`
	Status      string    `json:"status"`
}

var (
	tasks         []Task
	taskIDCounter int = 1
)


func toggleTaskStatus(status string) {
	row, _ := table.GetSelection()
	if row > 0 && row <= len(tasks) {
		tasks[row-1].Status = status
		refreshTable()
	}
}


func addTaskForm() {
	form := tview.NewForm().
		AddInputField("Description", "", 30, nil, nil).
		AddInputField("Due Date (YYYY-MM-DD)", "", 10, nil, nil).
		AddDropDown("Priority", []string{"🔥 High", "👍 Medium", "⭐ Low"}, 0, nil)

	form.AddButton("Save", func() {
		desc := form.GetFormItem(0).(*tview.InputField).GetText()
		dueDateStr := form.GetFormItem(1).(*tview.InputField).GetText()

		priority, _ := form.GetFormItem(2).(*tview.DropDown).GetCurrentOption()

		// Parse the due date
		dueDate, err := time.Parse("2006-01-02", dueDateStr)
		if err != nil {
			showMessage("Invalid date format. Use YYYY-MM-DD.")
			return
		}

		// Add the new task
		newTask := Task{
			ID:          taskIDCounter,
			Description: desc,
			Due:         dueDate,
			Priority:    []string{"🔥 High", "👍 Medium", "⭐ Low"}[priority],
			Status:      "❌ Pending",
		}
		tasks = append(tasks, newTask)
		taskIDCounter++

		// Refresh the table and return to the main view
		refreshTable()
		app.SetRoot(table, true)
	})

	form.AddButton("Cancel", func() {
		app.SetRoot(table, true)
	})

	// Set the form as the root view
	app.SetRoot(form, true)
}



func confirmDeleteTask() {
	row, _ := table.GetSelection()
	if row <= 0 || row > len(tasks) {
		showMessage("No task selected.")
		return
	}

	modal := tview.NewModal().
		SetText("Are you sure you want to delete this task?").
		AddButtons([]string{"Yes", "No"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == "Yes" {
				tasks = append(tasks[:row-1], tasks[row:]...)
				refreshTable()
			}
			app.SetRoot(table, true)
		})

	app.SetRoot(modal, true)
}


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





func showHelp() {
	modal := tview.NewModal().SetText("A: Add | E: Edit | D: Delete | C: Complete | P: Pending | S: Save | H: Help | Q: Quit").AddButtons([]string{"OK"}).SetDoneFunc(func(int, string) {
		app.SetRoot(table, true).SetFocus(table)
	})
		
	app.SetRoot(modal, true)
}



func showError(message string) {
	modal := tview.NewModal().SetText(message).AddButtons([]string{"OK"}).SetDoneFunc(func(int, string) {
		app.SetRoot(table, true).SetFocus(table)
	})
	app.SetRoot(modal, true)
}

func setupHeaders() {
	table.SetCell(0, 0, tview.NewTableCell("ID").SetAlign(tview.AlignCenter).SetSelectable(false))
	table.SetCell(0, 1, tview.NewTableCell("Description").SetAlign(tview.AlignCenter).SetSelectable(false))
	table.SetCell(0, 2, tview.NewTableCell("Due Date").SetAlign(tview.AlignCenter).SetSelectable(false))
	table.SetCell(0, 3, tview.NewTableCell("Priority").SetAlign(tview.AlignCenter).SetSelectable(false))
	table.SetCell(0, 4, tview.NewTableCell("Status").SetAlign(tview.AlignCenter).SetSelectable(false))
}

func refreshTable() {
	table.Clear()
	setupHeaders()
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