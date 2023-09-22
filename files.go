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

	progressCh := make(chan filewalker.ProgressInfo)

	go func(p <-chan filewalker.ProgressInfo) {
		for x := range p {
			progressLabel.SetText(fmt.Sprintf("Found %d images out of %d files", x.ImageFiles, x.TotalFiles))
		}
	}(progressCh)

	database.SetImageBasePath(dir)
	allImageFiles := filewalker.WalkForMe(dir, progressCh)
	database.InsertImagePathsIntoDatabase(allImageFiles)
	progress.Stop()
	progressWindow.Close()
}
