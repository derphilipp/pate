package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func swipeHandler(app fyne.App) {
	swipeWindow := app.NewWindow("Swipe Images")
	swipeWindow.Resize(fyne.NewSize(640, 480)) // Set the window size to 640x480

	// Load the first undecided image
	imagePath, err := getUndecidedImage()
	if err != nil {
		dialog.ShowError(err, swipeWindow)
		return
	}

	image := canvas.NewImageFromFile(imagePath)
	image.FillMode = canvas.ImageFillContain // Set FillMode to ImageFillContain

	// Check if the image resource is valid
	if image.Resource == nil {
		updateDecision(imagePath, "broken")
		refreshSwipeWindow(swipeWindow, app)
		return
	}

	// Set the image size to match the window size
	image.SetMinSize(fyne.NewSize(640, 480))

	leftBtn := widget.NewButton("Left", func() {
		updateDecision(imagePath, "not_copied")
		refreshSwipeWindow(swipeWindow, app)
	})

	rightBtn := widget.NewButton("Right", func() {
		updateDecision(imagePath, "copied")
		refreshSwipeWindow(swipeWindow, app)
	})

	// Handle left and right arrow key presses
	swipeWindow.Canvas().SetOnTypedKey(func(ev *fyne.KeyEvent) {
		switch ev.Name {
		case fyne.KeyLeft:
			leftBtn.Tapped(nil)
		case fyne.KeyRight:
			rightBtn.Tapped(nil)
		}
	})

	// Create a horizontal box layout with the left button, image, and right button
	content := container.New(layout.NewHBoxLayout(), leftBtn, image, rightBtn)
	swipeWindow.SetContent(content)
	swipeWindow.Show()
}

func refreshSwipeWindow(win fyne.Window, app fyne.App) {
	win.Close()
	swipeHandler(app)
}