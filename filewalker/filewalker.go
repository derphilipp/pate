package filewalker

import (
	"io/fs"
	"path/filepath"
	"strings"
	"sync"

	"github.com/charlievieth/fastwalk"
)

var conf = fastwalk.Config{
	Follow: false,
}

// this defines a structure of two interger: one for the current amount already processed and one for the total amount to process
type ProgressInfo struct {
	ImageFiles int
	TotalFiles int
}

func WalkForMe(dir string, progressCh chan<- ProgressInfo) []string {
	var totalFiles, imageFiles int
	var allImageFiles []string

	var mu sync.Mutex

	// Create a channel to control the number of concurrent goroutines
	sem := make(chan struct{}, 16) // Assuming 4 cores, adjust as needed

	var wg sync.WaitGroup

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
			progressCh <- ProgressInfo{imageFiles, totalFiles}
			mu.Unlock()
		}(path)
		return nil
	}

	fastwalk.Walk(&conf, dir, walkFn)
	return allImageFiles
}

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
