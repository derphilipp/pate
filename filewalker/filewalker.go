package filewalker

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"
	"time"

	"github.com/charlievieth/fastwalk"
	"github.com/sirupsen/logrus"
)

var conf = fastwalk.Config{
	Follow: false,
}

type Progress struct {
	FoundFiles    int
	SearchedFiles int
}

func SearchImageFiles(root string, fileCh chan<- string, progressCh chan<- Progress) {
	defer close(fileCh)
	defer close(progressCh)

	var totalSearched int
	var totalFound int

	ticker := time.NewTicker(10 * time.Millisecond)
	defer ticker.Stop()

	go func() {
		for range ticker.C {
			progressCh <- Progress{
				FoundFiles:    totalFound,
				SearchedFiles: totalSearched,
			}
		}
	}()

	walkFn := func(path string, info fs.DirEntry, err error) error {
		if err != nil {
			logrus.Warnf("Could not access directory: %v", err)
			return nil
		}

		if isImage(path) {
			fileCh <- path
			totalFound++
		}
		totalSearched++
		return nil
	}

	err := fastwalk.Walk(&conf, root, walkFn)
	if err != nil {
		fmt.Printf("Error walking directory: %v\n", err)
	}

	// One last report to the progress channel
	if progressCh != nil {
		progressCh <- Progress{
			FoundFiles:    totalFound,
			SearchedFiles: totalSearched,
		}
	}
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
