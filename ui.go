package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var (
	app   *tview.Application
	table *tview.Table
)

func initializeUI() error {
	// Create a new Tview application
	app = tview.NewApplication()

	// Create a table with borders
	table = tview.NewTable().SetBorders(true)

	// Set up the header row
	setupHeaders()

	// Populate the table with task data
	refreshTable()

	// Make the table selectable and set up keybindings
	table.SetSelectable(true, false).SetInputCapture(handleKeypress)

	// Create a header text view for the agenda title
	title := tview.NewTextView().
		SetText("📝 My Terminal To-Do Agenda").
		SetTextAlign(tview.AlignCenter).
		SetTextColor(tcell.ColorGreen)

	// Create a footer text view for the help information
	help := tview.NewTextView().
		SetText("A: Add | E: Edit | D: Delete | C: Complete | P: Pending | S: Save | H: Help | Q: Quit").
		SetTextAlign(tview.AlignCenter).
		SetTextColor(tcell.ColorYellow)

	// Arrange the title, table, and help in a vertical layout
	layout := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(title, 1, 1, false).
		AddItem(table, 0, 10, true).
		AddItem(help, 1, 1, false)

	// Run the application
	return app.SetRoot(layout, true).Run()
}
