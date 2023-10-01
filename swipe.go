package main

import (
	"fmt"
	"image"
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/derphilipp/pate/database"
	"github.com/sirupsen/logrus"
)

func swipeHandler(app fyne.App) {
	var loadImage func()
	swipeWindow := app.NewWindow("Swipe Images")
	swipeWindow.Resize(fyne.NewSize(640, 480)) // Set the window size to 640x480
	imageDataSource := alwaysHaveNImagesLoaded(15)
	// Function to load and display an image
	loadImage = func() {
		// Load the first undecided image
		// imagePath, err := database.GetUndecidedImage()

		// if err != nil {
		// 	// Check if there are no more undecided images
		// 	if err == sql.ErrNoRows {
		// 		dialog.ShowInformation("All Done", "No more images to display.", swipeWindow)
		// 		swipeWindow.Close()
		// 		return
		// 	}
		// 	dialog.ShowError(err, swipeWindow)
		// 	return
		// }
		imageToWorkOn := <-imageDataSource
		imageInFyne := canvas.NewImageFromImage(imageToWorkOn.image)
		imageInFyne.FillMode = canvas.ImageFillContain // Set FillMode to ImageFillContain

		// Check if the image resource is valid
		// if image.Resource == nil {
		// 	fmt.Printf("BRK")
		// 	updateDecision(imagePath, "broken")
		// 	loadImage() // Load the next image
		// 	return
		// }

		// Set the image size to match the window size
		imageInFyne.SetMinSize(fyne.NewSize(640, 480))

		leftBtn := widget.NewButton("Left", func() {
			go database.UpdateDecision(imageToWorkOn.path, "not_copied")
			loadImage() // Load the next image
		})

		rightBtn := widget.NewButton("Right", func() {
			go database.UpdateDecision(imageToWorkOn.path, "copied")
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
		imageNameLabel := widget.NewLabel(filepath.Base(imageToWorkOn.path))

		// Create a horizontal box layout with the left button, image, and right button
		content := container.NewVBox(
			container.New(layout.NewHBoxLayout(), leftBtn, imageInFyne, rightBtn),
			imageNameLabel,
		)
		swipeWindow.SetContent(content)
	}

	// Initially load the first image
	loadImage()

	swipeWindow.Show()
}

// func refreshSwipeWindow(win fyne.Window, app fyne.App) {
// 	win.Close()
// 	swipeHandler(app)
// }

type ImageCache struct {
	image image.Image
	path  string
}

func alwaysHaveNImagesLoaded(size int) <-chan ImageCache {
	imageCh := make(chan ImageCache, size)
	go func(imgageCh chan<- ImageCache) {
		allImagePaths, err := database.GetAllUndecidedPaths()
		if err != nil {
			logrus.Fatal("Failed to load images:", err)
			return
		}
		for _, img := range allImagePaths {
			imagePath := img
			imageData := database.LoadSingleFile(imagePath)
			imageCh <- ImageCache{image: imageData, path: imagePath}
			fmt.Printf("Image cached: %s\n", imagePath)
		}
	}(imageCh)
	return imageCh
}

// func loadNextImage() string {
// 	if len(imageCache) == 0 {
// 		// If the cache is empty, preload the next N images
// 		paths, err := database.GetNextNImages(5) // Preload the next 5 images
// 		if err != nil {
// 			log.Println("Failed to load images:", err)
// 			return ""
// 		}
// 		imageCache = append(imageCache, paths...)
// 	}

// 	// Get the next image from the cache
// 	imagePath := imageCache[0]
// 	imageCache = imageCache[1:]

// 	return imagePath
// }
