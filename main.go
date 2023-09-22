package main

import (
	"fmt"
	"sync"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/derphilipp/pate/checksum"
	"github.com/derphilipp/pate/database"
)

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Fyne App")

	// Initialize SQLite database
	database.InitDatabase()

	// Load existing files from the database
	database.LoadExistingFiles()

	// Button handlers
	selectInputHandler := func() {
		tempWindow := myApp.NewWindow("Open Folder")
		tempWindow.Resize(fyne.NewSize(1024, 786))
		tempWindow.Show()

		fmt.Printf("Current time and date before input: %s\n", time.Now().Format("2006-01-02 15:04:05"))

		dialog.ShowFolderOpen(func(dir fyne.ListableURI, err error) {
			if err == nil && dir != nil {
				go scanDirectory(dir.Path(), myApp)
			}
			tempWindow.Close()
			fmt.Printf("Current time and date after input: %s\n", time.Now().Format("2006-01-02 15:04:05"))
		}, tempWindow)
	}

	calculateChecksumsHandler := func() {
		go ChecksumFiles(myApp)
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

// func CalcAllChecksum(){

// }

func ChecksumFiles(app fyne.App) {
	fmt.Printf("Current time and date before checksum: %s\n", time.Now().Format("2006-01-02 15:04:05"))

	amount := database.CountNonchecksummedFiles()
	uncheckedImages, _ := database.GetUnchecksummedImagesFromDatabase()
	progress := widget.NewProgressBar()
	progressLabel := widget.NewLabel("Checksumming files")
	progressContainer := container.NewVBox(progressLabel, progress)
	progressWindow := app.NewWindow("Checksum Progress")

	progressWindow.SetContent(progressContainer)
	progressWindow.Show()

	const numWorkers = 16

	inputChan := make(chan string, amount)
	progressCh := make(chan int, numWorkers)
	batchChan := make(chan database.FileChecksum, amount)

	go database.ChecksumWriter(batchChan, nil)

	var wg sync.WaitGroup

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go checksum.ChecksumWorker(inputChan, batchChan, progressCh, &wg)
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
		close(progressCh)
		close(batchChan)
	}()

	processedCount := 0

	for justProcessed := range progressCh {
		processedCount += justProcessed
		progress.SetValue(float64(processedCount) / float64(amount))
		progressLabel.SetText(fmt.Sprintf("Processed %d out of %d files", processedCount, amount))

		if int64(processedCount) == amount {
			break
		}
	}

	progressWindow.Close()
	fmt.Printf("Current time and date after checksum: %s\n", time.Now().Format("2006-01-02 15:04:05"))
}
