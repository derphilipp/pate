package validimage

import (
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"path/filepath"
	"strings"

	_ "golang.org/x/image/bmp"
	_ "golang.org/x/image/tiff"
	_ "golang.org/x/image/webp"
)

// IsValidImage checks if the provided file is a valid image.
func IsValidImage(filePath string) (bool, error) {
	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		return false, err
	}
	defer file.Close()

	// Check the file extension
	ext := strings.ToLower(filepath.Ext(filePath))
	validExtensions := []string{".jpg", ".jpeg", ".png", ".gif", ".bmp", ".tiff", ".webp"}
	isValidExt := false
	for _, v := range validExtensions {
		if ext == v {
			isValidExt = true
			break
		}
	}
	if !isValidExt {
		return false, fmt.Errorf("unsupported file extension")
	}

	// Decode the image
	_, _, err = image.Decode(file)
	if err != nil {
		return false, err
	}

	return true, nil
}
