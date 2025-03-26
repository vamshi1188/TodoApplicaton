package main

import "github.com/gdamore/tcell/v2"

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
		saveTasksToS3()
		showMessage("Tasks saved successfully!")
	case 'h', 'H':
		showHelp()
	case 'q', 'Q':
		app.Stop()
	}
	return event
}
