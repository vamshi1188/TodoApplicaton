package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
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

var (
	tasks         []Task
	taskIDCounter int = 1
	app           *tview.Application
	table         *tview.Table
	bucketName    = "vamshigodev" // Updated with your S3 bucket name
	objectKey     = "tasks.json"
	s3Client      *s3.Client
	uploader      *manager.Uploader
	downloader    *manager.Downloader
)

func main() {
	// Initialize AWS S3 client
	initS3Client()

	// Load tasks from S3
	loadTasksFromS3()

	// Create a new Tview application
	app = tview.NewApplication()

	// Create a table with borders
	table = tview.NewTable().SetBorders(true)

<<<<<<< HEAD
	// Set up the header row
=======
	// Set up the header row.
>>>>>>> 5a028547f7b5111145029132d4fb3f14a9d658cb
	setupHeaders()

	// Populate the table with task data
	refreshTable()

<<<<<<< HEAD
	// Make the table selectable and set up keybindings
	table.SetSelectable(true, false).SetInputCapture(handleKeypress)

	// Create a header text view for the agenda title
=======
	// Make the table selectable and set up keybindings.
	table.SetSelectable(true, false).SetInputCapture(handleKeypress)

	// Create a header text view for the agenda title.
>>>>>>> 5a028547f7b5111145029132d4fb3f14a9d658cb
	title := tview.NewTextView().
		SetText("📝 My Terminal To-Do Agenda").
		SetTextAlign(tview.AlignCenter).
		SetTextColor(tcell.ColorGreen)

<<<<<<< HEAD
	// Create a footer text view for the help information
	help := tview.NewTextView().
		SetText("A: Add | E: Edit | D: Delete | C: Complete | P: Pending | S: Save | H: Help | Q: Quit").
		SetTextAlign(tview.AlignCenter).
		SetTextColor(tcell.ColorYellow)
=======
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
>>>>>>> 5a028547f7b5111145029132d4fb3f14a9d658cb

	// Arrange the title, table, and help in a vertical layout
	layout := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(title, 1, 1, false). // Title at the top
		AddItem(table, 0, 10, true). // Table in the middle
		AddItem(help, 1, 1, false)   // Help at the bottom

	// Run the application
	if err := app.SetRoot(layout, true).Run(); err != nil {
		log.Fatalf("Error running application: %v", err)
	}
}

<<<<<<< HEAD
// Initialize the AWS S3 client
func initS3Client() {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("ap-south-1"))
	if err != nil {
		log.Fatalf("Unable to load AWS SDK config: %v", err)
	}

	s3Client = s3.NewFromConfig(cfg)
	uploader = manager.NewUploader(s3Client)
	downloader = manager.NewDownloader(s3Client)
}

// Load tasks from S3
func loadTasksFromS3() {
	buf := manager.NewWriteAtBuffer([]byte{})
	_, err := downloader.Download(context.TODO(), buf, &s3.GetObjectInput{
		Bucket: &bucketName,
		Key:    &objectKey,
	})
	if err != nil {
		log.Printf("Failed to download tasks from S3: %v", err)
		return
	}

	if err := json.Unmarshal(buf.Bytes(), &tasks); err != nil {
		log.Printf("Failed to unmarshal tasks: %v", err)
		return
	}

	// Update taskIDCounter to avoid ID conflicts
	for _, task := range tasks {
		if task.ID >= taskIDCounter {
			taskIDCounter = task.ID + 1
		}
	}
}

// Set up the table headers
=======
// Set up the table headers.
>>>>>>> 5a028547f7b5111145029132d4fb3f14a9d658cb
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

<<<<<<< HEAD
// Refresh the table with updated task data
=======
// Refresh the table with updated task data.
>>>>>>> 5a028547f7b5111145029132d4fb3f14a9d658cb
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

<<<<<<< HEAD
// Handle keypress events
=======
// Handle keypress events.
>>>>>>> 5a028547f7b5111145029132d4fb3f14a9d658cb
func handleKeypress(event *tcell.EventKey) *tcell.EventKey {
	switch event.Rune() {
	case 'c', 'C':
		toggleTaskStatus("✅ Done")
<<<<<<< HEAD
	case 'p', 'P':
=======
	case 'x', 'X':
>>>>>>> 5a028547f7b5111145029132d4fb3f14a9d658cb
		toggleTaskStatus("❌ Pending")
	case 'a', 'A':
		addTaskForm()
	case 'd', 'D':
		confirmDeleteTask()
	case 'e', 'E':
		editTaskForm()
	case 's', 'S':
<<<<<<< HEAD
		saveTasksToS3()
=======
		saveTasks()
>>>>>>> 5a028547f7b5111145029132d4fb3f14a9d658cb
		showMessage("Tasks saved successfully!")
	case 'h', 'H':
		showHelp()
	case 'q', 'Q':
		app.Stop()
	}
	return event
}

<<<<<<< HEAD
// Toggle task status between "✅ Done" and "❌ Pending"
=======
// Toggle task status between "✅ Done" and "❌ Pending".
>>>>>>> 5a028547f7b5111145029132d4fb3f14a9d658cb
func toggleTaskStatus(status string) {
	row, _ := table.GetSelection()
	if row > 0 && row <= len(tasks) {
		tasks[row-1].Status = status
		refreshTable()
	}
}

// Add a new task using a form
func addTaskForm() {
	form := tview.NewForm().
		AddInputField("Description", "", 30, nil, nil).
		AddInputField("Due Date (YYYY-MM-DD)", "", 10, nil, nil).
		AddDropDown("Priority", []string{"🔥 High", "👍 Medium", "⭐ Low"}, 0, nil)

	form.AddButton("Save", func() {
		desc := form.GetFormItem(0).(*tview.InputField).GetText()
		dueDateStr := form.GetFormItem(1).(*tview.InputField).GetText()
<<<<<<< HEAD
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

// Edit an existing task using a form
func editTaskForm() {
	row, _ := table.GetSelection()
	if row <= 0 || row > len(tasks) {
		showMessage("No task selected.")
		return
	}

	task := tasks[row-1]
	form := tview.NewForm().
		AddInputField("Description", task.Description, 30, nil, nil).
		AddInputField("Due Date (YYYY-MM-DD)", task.Due.Format("2006-01-02"), 10, nil, nil).
		AddDropDown("Priority", []string{"🔥 High", "👍 Medium", "⭐ Low"}, getPriorityIndex(task.Priority), nil)

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

		// Update the task
		task.Description = desc
		task.Due = dueDate
		task.Priority = []string{"🔥 High", "👍 Medium", "⭐ Low"}[priority]

		refreshTable()
		app.SetRoot(table, true)
	})

	form.AddButton("Cancel", func() {
		app.SetRoot(table, true)
	})

	app.SetRoot(form, true)
}

// Get the index of the priority in the dropdown
func getPriorityIndex(priority string) int {
	switch priority {
	case "🔥 High":
		return 0
	case "👍 Medium":
		return 1
	case "⭐ Low":
		return 2
	default:
		return 0
	}
}

// Save tasks to S3
func saveTasksToS3() {
	data, err := json.Marshal(tasks)
	if err != nil {
		log.Printf("Failed to marshal tasks: %v", err)
		showMessage("Failed to save tasks.")
		return
	}

	_, err = uploader.Upload(context.TODO(), &s3.PutObjectInput{
		Bucket:      &bucketName,
		Key:         &objectKey,
		Body:        bytes.NewReader(data),
		ContentType: aws.String("application/json"), // Set Content-Type to application/json
	})
	if err != nil {
		log.Printf("Failed to upload tasks to S3: %v", err)
		showMessage("Failed to save tasks.")
		return
	}

	showMessage("Tasks saved successfully!")
}

// Confirm delete task
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

// Show help information in a modal dialog
func showHelp() {
	helpText := `Help:
- A: Add a new task
- E: Edit the selected task
- D: Delete the selected task
- C: Mark the selected task as complete
- P: Mark the selected task as pending
- S: Save tasks to S3
- H: Show this help
- Q: Quit the application`

	modal := tview.NewModal().
		SetText(helpText).
		AddButtons([]string{"OK"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			app.SetRoot(table, true)
		})

=======
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
>>>>>>> 5a028547f7b5111145029132d4fb3f14a9d658cb
	app.SetRoot(modal, true)
}

// Show a message in a modal dialog
func showMessage(message string) {
	modal := tview.NewModal().
		SetText(message).
		AddButtons([]string{"OK"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			app.SetRoot(table, true)
		})

	app.SetRoot(modal, true)
}
