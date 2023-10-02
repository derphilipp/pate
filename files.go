package main

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/derphilipp/pate/database"
	"github.com/derphilipp/pate/filewalker"
)

func scanDirectory(dir string, app fyne.App) {
	// findFilesProgres := widget.NewProgressBarInfinite()
	progress := widget.NewProgressBarInfinite()
	progressLabel := widget.NewLabel("Searching files")
	progressContainer := container.NewVBox(progressLabel, progress)
	progressWindow := app.NewWindow("Scanning Progress")
	progressWindow.SetContent(progressContainer)
	progressWindow.Show()

	progressCh := make(chan filewalker.Progress, 100)
	fileWalkerCh := make(chan string, 1024)

	go func(p <-chan filewalker.Progress) {
		for x := range p {
			progressLabel.SetText(fmt.Sprintf("Found %d images out of %d files", x.FoundFiles, x.SearchedFiles))
		}
	}(progressCh)

	database.SetImageBasePath(dir)

	go filewalker.SearchImageFiles(dir, fileWalkerCh, progressCh)
	var allImageFiles []string
	for image := range fileWalkerCh {
		allImageFiles = append(allImageFiles, image)
	}
	fmt.Printf("\nAll done, found %d images\n", len(allImageFiles))
	// allImageFiles := filewalker.WalkForMe(dir, progressCh)

	database.InsertImagePathsIntoDatabase(allImageFiles)
	progress.Stop()
	progressWindow.Close()
}
