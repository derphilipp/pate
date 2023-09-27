package main

import (
	"fmt"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/derphilipp/pate/checksum"
	"github.com/derphilipp/pate/database"
	"github.com/sirupsen/logrus"

	"github.com/sourcegraph/conc/pool"
)

type ChecksumProgress struct {
	processedCount int64
	amount         int64
}

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

	findDuplicatesHandler := func() {
		go DuplicateFiles(myApp)
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
	findDuplicatesBtn := widget.NewButton("Find Duplicates", findDuplicatesHandler)
	swipeBtn := widget.NewButton("Swipe", swipeHandlerFunc)
	selectOutputBtn := widget.NewButton("Select Output", selectOutputHandler)
	copyBtn := widget.NewButton("Copy", copyHandler)
	exitBtn := widget.NewButton("Exit", exitHandler)

	// Add buttons to a vertical box layout
	content := container.NewVBox(
		selectInputBtn,
		calculateChecksumsBtn,
		findDuplicatesBtn,
		swipeBtn,
		selectOutputBtn,
		copyBtn,
		exitBtn,
	)

	myWindow.SetContent(content)
	myWindow.ShowAndRun()
}

func CalcAllChecksum(totalProgressCh chan<- ChecksumProgress) {
	const numWorkers = 24

	amount := database.CountNonchecksummedFiles()
	uncheckedImages, _ := database.GetUnchecksummedImagesFromDatabase()

	batchChan := make(chan database.FileChecksum, 256)

	go database.ChecksumWriter(batchChan, nil)

	p := pool.New().WithMaxGoroutines(numWorkers)
	for _, uncheckedImage := range uncheckedImages {
		uncheckedImage := uncheckedImage
		p.Go(func() {
			sum, err := checksum.CalculateChecksum(uncheckedImage)
			if err != nil {
				logrus.Warnf("Error calculating checksum for %s: %v\n", uncheckedImage, err)
			} else {
				// fmt.Printf("Checksum for %s: %s\n", uncheckedImage, sum)
				batchChan <- database.FileChecksum{FilePath: uncheckedImage, Checksum: sum}
			}
			totalProgressCh <- ChecksumProgress{processedCount: 1, amount: amount}
			//		progressCh <- 1
		})
	}
	p.Wait()
	close(batchChan)
}

func DuplicateFiles(app fyne.App) {
	database.DetectAndHandleDuplicates()
}

func ChecksumFiles(app fyne.App) {
	fmt.Printf("Current time and date before checksum: %s\n", time.Now().Format("2006-01-02 15:04:05"))

	progress := widget.NewProgressBar()
	progressLabel := widget.NewLabel("Checksumming files")
	progressContainer := container.NewVBox(progressLabel, progress)
	progressWindow := app.NewWindow("Checksum Progress")

	progressWindow.SetContent(progressContainer)
	progressWindow.Show()
	checksumProgressChan := make(chan ChecksumProgress, 1024)

	go func(progressChan <-chan ChecksumProgress) {
		var processedCount int64 = 0
		for justProcessed := range progressChan {
			processedCount += justProcessed.processedCount
			progress.SetValue(float64(processedCount) / float64(justProcessed.amount))
			// progressCh <- progress.SetValue(float64(processedCount) / float64(amount))
			labelText := fmt.Sprintf("Processed %d out of %d files", processedCount, justProcessed.amount)
			// fmt.Println(labelText)
			progressLabel.SetText(labelText)

		}
	}(checksumProgressChan)

	CalcAllChecksum(checksumProgressChan)
	// xxxx

	progressWindow.Close()
	fmt.Printf("Current time and date after checksum: %s\n", time.Now().Format("2006-01-02 15:04:05"))
}
