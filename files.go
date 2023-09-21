package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/charlievieth/fastwalk"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// Supported image extensions
var imageExtensions = []string{".jpg", ".jpeg", ".png", ".gif", ".bmp", ".tiff", ".webp"}

func isImage(path string) bool {
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

func scanDirectory(dir string, app fyne.App) {
	// findFilesProgres := widget.NewProgressBarInfinite()
	progress := widget.NewProgressBarInfinite()
	progressLabel := widget.NewLabel("Searching files")
	progressContainer := container.NewVBox(progressLabel, progress)
	progressWindow := app.NewWindow("Scanning Progress")
	progressWindow.SetContent(progressContainer)
	progressWindow.Show()

	var totalFiles, imageFiles int
	var allImageFiles []string

	var mu sync.Mutex

	// Create a channel to control the number of concurrent goroutines
	sem := make(chan struct{}, 4) // Assuming 4 cores, adjust as needed

	var wg sync.WaitGroup

	conf := fastwalk.Config{
		Follow: false,
	}

	walkFn := func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}

		sem <- struct{}{} // Acquire a token
		wg.Add(1)

		go func(p string) {
			defer func() {
				<-sem // Release the token
				wg.Done()
			}()

			if isImage(p) {
				mu.Lock()
				imageFiles++
				allImageFiles = append(allImageFiles, p)
				mu.Unlock()
			}

			mu.Lock()
			totalFiles++
			if totalFiles%10 == 0 {
				progressLabel.SetText(fmt.Sprintf("Found %d images out of %d files", imageFiles, totalFiles))
			}
			mu.Unlock()
		}(path)

		return nil
	}

	fastwalk.Walk(&conf, dir, walkFn)

	wg.Wait() // Wait for all goroutines to finish

	insertImagePathsIntoDatabase(allImageFiles)
	progress.Stop()
	fmt.Printf("DONE")
	progressWindow.Close()
}

func checksumFiles(app fyne.App) {
	amount := countNonchecksummedFiles()
	uncheckedImages, _ := getUnchecksummedImagesFromDatabase()
	progress := widget.NewProgressBar()
	progressLabel := widget.NewLabel("Checksumming files")
	progressContainer := container.NewVBox(progressLabel, progress)
	progressWindow := app.NewWindow("Checksum Progress")

	progressWindow.SetContent(progressContainer)
	progressWindow.Show()

	const numWorkers = 16

	inputChan := make(chan string, amount)
	doneChan := make(chan int, numWorkers)
	batchChan := make(chan FileChecksum, amount)

	go dbWriter(batchChan)

	var wg sync.WaitGroup

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go checksumWorker(inputChan, batchChan, doneChan, &wg)
	}

	go func() {
		for _, filePath := range uncheckedImages {
			inputChan <- filePath
		}
		fmt.Printf("Done inputting....")
		close(inputChan)
	}()

	go func() {
		wg.Wait()
		close(doneChan)
		close(batchChan)
	}()

	processedCount := 0
	for {
		select {
		// case result := <-outputChan:
		//	fmt.Printf("File: %s, Checksum: %s\n", result.FilePath, result.FilePath)

		case justProcessed := <-doneChan:
			// fmt.Printf("Just done %d\n", justProcessed)
			processedCount += justProcessed
			progress.SetValue(float64(processedCount) / float64(amount))
			progressLabel.SetText(fmt.Sprintf("Processed %d out of %d files", processedCount, amount))
			// fmt.Printf("Processed %d out of %d files\n", processedCount, amount)

			if int64(processedCount) == amount {
				break
			}
		}
	}
	/*
		for i, uncheckedImage := range uncheckedImages {
			chk, _ := calculateChecksum(uncheckedImage)
			updateChecksumInDatabase(uncheckedImage, chk)
			progressLabel.SetText(fmt.Sprintf("Calculated %d out of %d checksums", i, amount))
			progress.SetValue(float64(i) / float64(amount))
		}
	*/
	progressWindow.Close()
}

// func checksumFilesPlain(app fyne.App) {
// 	amount := countNonchecksummedFiles()
// 	uncheckedImages, _ := getUnchecksummedImagesFromDatabase()
// 	progress := widget.NewProgressBar()
// 	progressLabel := widget.NewLabel("Checksumming files")
// 	progressContainer := container.NewVBox(progressLabel, progress)
// 	progressWindow := app.NewWindow("Checksum Progress")

// 	progressWindow.SetContent(progressContainer)
// 	progressWindow.Show()

// 	for i, uncheckedImage := range uncheckedImages {
// 		chk, _ := calculateChecksum(uncheckedImage)
// 		updateChecksumInDatabase(uncheckedImage, chk)
// 		progressLabel.SetText(fmt.Sprintf("Calculated %d out of %d checksums", i, amount))
// 		progress.SetValue(float64(i) / float64(amount))
// 	}

// 	progressWindow.Close()
// }

func calculateChecksum(filePath string) (string, error) {
	fileData, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	hash := sha256.Sum256(fileData)
	return hex.EncodeToString(hash[:]), nil
}

// func calculateChecksums(app fyne.App) {
// 	progressWindow := app.NewWindow("Calculating Checksums")
// 	progressWindow.Resize(fyne.NewSize(400, 150))

// 	// Create a progress bar
// 	progressBar := widget.NewProgressBarInfinite()

// 	// Label to display number of images found
// 	imageCountLabel := widget.NewLabel("Images found: 0")

// 	content := container.NewVBox(
// 		widget.NewLabel("Calculating checksums..."),
// 		progressBar,
// 		imageCountLabel,
// 	)
// 	progressWindow.SetContent(content)
// 	progressWindow.Show()

// 	go func() {
// 		// Assuming you have a function getImagesFromDatabase() that retrieves all images from the database
// 		images, err := getImagesFromDatabase()
// 		if err != nil {
// 			log.Println("Error getting images from database:", err)
// 			return
// 		}

// 		for i, imagePath := range images {
// 			// Assuming you have a function calculateChecksum(imagePath) that calculates the checksum of an image
// 			checksum, err := calculateChecksum(imagePath)
// 			if err != nil {
// 				log.Println("Error calculating checksum:", err)
// 				continue
// 			}

// 			// Update the database with the checksum
// 			updateChecksumInDatabase(imagePath, checksum)

// 			// Update the label with the number of images processed
// 			fyne.CurrentApp().Driver().RunOnMain(func() {
// 				imageCountLabel.SetText(fmt.Sprintf("Images found: %d", i+1))
// 			})
// 		}

// 		// Close the progress window when done
// 		fyne.CurrentApp().Driver().
// 			fyne.CurrentApp().Driver().RunOnMain(func() {
// 			progressWindow.Close()
// 		})
// 	}()
// }

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
