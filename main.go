package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Fyne App")

	// Initialize SQLite database
	initDatabase()

	// Load existing files from the database
	loadExistingFiles()

	// Button handlers
	selectInputHandler := func() {
		dialog.ShowFolderOpen(func(dir fyne.ListableURI, err error) {
			if err == nil && dir != nil {
				go scanDirectory(dir.Path(), myApp)
			}
		}, myWindow)
	}

	calculateChecksumsHandler := func() {
		go checksumFiles(myApp)
	}

	selectOutputHandler := func() {
		outputWindow := myApp.NewWindow("Select Output")
		outputWindow.SetContent(widget.NewLabel("Select Output Window Content"))
		outputWindow.Show()
	}

	copyHandler := func() {
		copyWindow := myApp.NewWindow("Copy")
		copyWindow.SetContent(widget.NewLabel("Copy Window Content"))
		copyWindow.Show()
	}

	exitHandler := func() {
		myApp.Quit()
	}

	swipeHandlerFunc := func() {
		swipeHandler(myApp)
	}

	// Create buttons
	selectInputBtn := widget.NewButton("Select Input", selectInputHandler)
	calculateChecksumsBtn := widget.NewButton("Calculate Checksums", calculateChecksumsHandler)
	swipeBtn := widget.NewButton("Swipe", swipeHandlerFunc)
	selectOutputBtn := widget.NewButton("Select Output", selectOutputHandler)
	copyBtn := widget.NewButton("Copy", copyHandler)
	exitBtn := widget.NewButton("Exit", exitHandler)

	// Add buttons to a vertical box layout
	content := container.NewVBox(
		selectInputBtn,
		calculateChecksumsBtn,
		swipeBtn,
		selectOutputBtn,
		copyBtn,
		exitBtn,
	)

	myWindow.SetContent(content)
	myWindow.ShowAndRun()
}
