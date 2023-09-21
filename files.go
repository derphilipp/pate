package main

import (
	"crypto/sha256"
	"encoding/hex"
	"os"
	"path/filepath"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func scanDirectory(dir string, app fyne.App) {
	progress := widget.NewProgressBar()
	progressLabel := widget.NewLabel("Calculating Checksums")
	progressContainer := container.NewVBox(progressLabel, progress)
	progressWindow := app.NewWindow("Scanning Progress")
	progressWindow.SetContent(progressContainer)
	progressWindow.Show()

	imageExtensions := []string{".jpg", ".jpeg", ".png", ".gif", ".bmp", ".tiff", ".webp"}
	var totalFiles, imageFiles int

	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		totalFiles++
		ext := strings.ToLower(filepath.Ext(path))
		for _, imgExt := range imageExtensions {
			if ext == imgExt {
				imageFiles++
				checksum, _ := calculateChecksum(path)
				statement, _ := db.Prepare("INSERT INTO images (path, checksum) VALUES (?, ?)")
				statement.Exec(path, checksum)
				break
			}
		}
		progress.SetValue(float64(imageFiles) / float64(totalFiles))
		return nil
	})

	progressWindow.Close()
}

func calculateChecksum(filePath string) (string, error) {
	fileData, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	hash := sha256.Sum256(fileData)
	return hex.EncodeToString(hash[:]), nil
}
