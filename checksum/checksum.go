package checksum

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"sync"

	"github.com/derphilipp/pate/database"
)

func ChecksumWorker(inputChan <-chan string, batchChan chan<- database.FileChecksum, progressCh chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()

	for filePath := range inputChan {
		checksum, err := calculateChecksum(filePath)
		progressCh <- 1
		if err != nil {
			fmt.Printf("Error calculating checksum for %s: %v\n", filePath, err)
			continue
		}

		batchChan <- database.FileChecksum{FilePath: filePath, Checksum: checksum}
	}
}

func calculateChecksum(filePath string) (string, error) {
	fileData, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	hash := sha256.Sum256(fileData)
	return hex.EncodeToString(hash[:]), nil
}
