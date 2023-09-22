package main

import (
	"database/sql"
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/derphilipp/pate/database"
)

func swipeHandler(app fyne.App) {
	var loadImage func()
	swipeWindow := app.NewWindow("Swipe Images")
	swipeWindow.Resize(fyne.NewSize(640, 480)) // Set the window size to 640x480

	// Function to load and display an image
	loadImage = func() {
		// Load the first undecided image
		imagePath, err := database.GetUndecidedImage()
		if err != nil {
			// Check if there are no more undecided images
			if err == sql.ErrNoRows {
				dialog.ShowInformation("All Done", "No more images to display.", swipeWindow)
				swipeWindow.Close()
				return
			}
			dialog.ShowError(err, swipeWindow)
			return
		}

		image := canvas.NewImageFromFile(imagePath)
		image.FillMode = canvas.ImageFillContain // Set FillMode to ImageFillContain

		// Check if the image resource is valid
		// if image.Resource == nil {
		// 	fmt.Printf("BRK")
		// 	updateDecision(imagePath, "broken")
		// 	loadImage() // Load the next image
		// 	return
		// }

		// Set the image size to match the window size
		image.SetMinSize(fyne.NewSize(640, 480))

		leftBtn := widget.NewButton("Left", func() {
			database.UpdateDecision(imagePath, "not_copied")
			loadImage() // Load the next image
		})

		rightBtn := widget.NewButton("Right", func() {
			database.UpdateDecision(imagePath, "copied")
			loadImage() // Load the next image
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

		// Display the image name below the image
		imageNameLabel := widget.NewLabel(filepath.Base(imagePath))

		// Create a horizontal box layout with the left button, image, and right button
		content := container.NewVBox(
			container.New(layout.NewHBoxLayout(), leftBtn, image, rightBtn),
			imageNameLabel,
		)
		swipeWindow.SetContent(content)
	}

	// Initially load the first image
	loadImage()

	swipeWindow.Show()
}

func refreshSwipeWindow(win fyne.Window, app fyne.App) {
	win.Close()
	swipeHandler(app)
}
