package main

import (
	"fmt"
	"sync"
)

var batchSize = 10

func checksumWorker(inputChan <-chan string, batchChan chan<- FileChecksum, doneChan chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()

	for filePath := range inputChan {
		checksum, err := calculateChecksum(filePath)
		if err != nil {
			fmt.Printf("Error calculating checksum for %s: %v\n", filePath, err)
			continue
		}

		batchChan <- FileChecksum{FilePath: filePath, Checksum: checksum}
		doneChan <- 1

	}
}
