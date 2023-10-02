package filewalker

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSearchImageFiles(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "testSearchImageFiles")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create some dummy image files in the directory
	imageFiles := []string{"image1.jpg", "image2.jpeg", "image3.png", "notImage.txt"}
	for _, file := range imageFiles {
		filePath := filepath.Join(tempDir, file)
		if err := os.WriteFile(filePath, []byte("test content"), 0o644); err != nil {
			t.Fatalf("Failed to create dummy file: %v", err)
		}
	}

	// Channels for SearchImageFiles function
	fileCh := make(chan string, 1024)
	progressCh := make(chan Progress, 100)

	// Call the SearchImageFiles function
	go SearchImageFiles(tempDir, fileCh, progressCh)

	// Collect the results
	var foundFiles []string
	for file := range fileCh {
		foundFiles = append(foundFiles, file)
	}

	var searchedFilesMatch int
	var searchedFilesTotal int

	// time.Sleep(1000 * time.Millisecond)
	for progress := range progressCh {
		searchedFilesMatch = progress.FoundFiles
		searchedFilesTotal = progress.SearchedFiles
	}

	// Check if the correct number of image files were detected
	if len(foundFiles) != 3 {
		t.Errorf("Expected to find 3 image files, but found %d", len(foundFiles))
	}

	// 1 Directory, 4 Files
	if searchedFilesTotal != 5 {
		t.Errorf("Expected to find 5 searched files in total, but found %d", searchedFilesTotal)
	}

	// 3 Image files
	if searchedFilesMatch != 3 {
		t.Errorf("Expected to find 3 searched files that match, but found %d", searchedFilesTotal)
	}
}
