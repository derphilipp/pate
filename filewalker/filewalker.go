package filewalker

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"
	"sync"
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
	var wg sync.WaitGroup
	var mutex sync.Mutex

	var totalSearched int
	var totalFound int

	ticker := time.NewTicker(10 * time.Millisecond)

	go func() {
		for range ticker.C {
			mutex.Lock()
			if progressCh != nil {
				progressCh <- Progress{
					FoundFiles:    totalFound,
					SearchedFiles: totalSearched,
				}
			}
			mutex.Unlock()
		}
	}()

	walkFn := func(path string, info fs.DirEntry, err error) error {
		if err != nil {
			logrus.Warnf("Could not access directory: %v", err)
			return nil
			// return err
		}

		wg.Add(1)
		go func() {
			defer wg.Done()
			// extension := strings.ToLower(filepath.Ext(path))
			// convert extension to lowercase
			if isImage(path) {
				// if extension == ".jpeg" {
				fileCh <- path
				mutex.Lock()
				totalFound++
				mutex.Unlock()
			}

			mutex.Lock()
			totalSearched++
			mutex.Unlock()
		}()

		return nil
	}

	err := fastwalk.Walk(&conf, root, walkFn)
	if err != nil {
		fmt.Printf("Error walking directory: %v\n", err)
	}

	wg.Wait()
	close(fileCh)
	ticker.Stop()
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
