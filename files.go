package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// Supported image extensions
var imageExtensions = []string{".jpg", ".jpeg", ".png", ".gif", ".bmp", ".tiff", ".webp"}

func scanDirectory(dirPath string, app fyne.App) {
	// Initialize the SQLite database (assuming you have a function for this)
	initDatabase()

	// Recursive function to scan directories for image files
	var scan func(path string)
	scan = func(path string) {
		files, err := os.ReadDir(path)
		if err != nil {
			log.Println("Error reading directory:", err)
			return
		}

		for _, file := range files {
			fullPath := filepath.Join(path, file.Name())

			if file.IsDir() {
				scan(fullPath)
			} else if isImageFile(fullPath) { // Assuming you have a function to check if a file is an image
				// Insert the image path into the SQLite database
				insertImagePathIntoDatabase(fullPath)
			}
		}
	}

	// Start scanning from the selected directory
	scan(dirPath)

	// Calculate checksums for the images found
	calculateChecksums(app)
}

func calculateChecksum(filePath string) (string, error) {
	fileData, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	hash := sha256.Sum256(fileData)
	return hex.EncodeToString(hash[:]), nil
}

func calculateChecksums(app fyne.App) {
	progressWindow := app.NewWindow("Calculating Checksums")
	progressWindow.Resize(fyne.NewSize(400, 150))

	// Create a progress bar
	progressBar := widget.NewProgressBarInfinite()

	// Label to display number of images found
	imageCountLabel := widget.NewLabel("Images found: 0")

	content := container.NewVBox(
		widget.NewLabel("Calculating checksums..."),
		progressBar,
		imageCountLabel,
	)
	progressWindow.SetContent(content)
	progressWindow.Show()

	go func() {
		// Assuming you have a function getImagesFromDatabase() that retrieves all images from the database
		images, err := getImagesFromDatabase()
		if err != nil {
			log.Println("Error getting images from database:", err)
			return
		}

		for i, imagePath := range images {
			// Assuming you have a function calculateChecksum(imagePath) that calculates the checksum of an image
			checksum, err := calculateChecksum(imagePath)
			if err != nil {
				log.Println("Error calculating checksum:", err)
				continue
			}

			// Update the database with the checksum
			updateChecksumInDatabase(imagePath, checksum)

			// Update the label with the number of images processed
			app.QueueUpdate(func() {
				imageCountLabel.SetText(fmt.Sprintf("Images found: %d", i+1))
			})
		}

		// Close the progress window when done
		app.QueueUpdate(func() {
			progressWindow.Close()
		})
	}()
}

func isImageFile(path string) bool {
	// Get the file extension
	ext := strings.ToLower(filepath.Ext(path))

	// Check if the extension is in the list of supported image extensions
	for _, imageExt := range imageExtensions {
		if ext == imageExt {
			return true
		}
	}
	return false
}
