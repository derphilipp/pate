package filewalker

import (
	"fmt"
	"io/fs"
	"sync"
	"time"

	"github.com/charlievieth/fastwalk"
)

type Progress struct {
	FoundFiles    int
	SearchedFiles int
}

func SearchJPEGFiles(root string, fileCh chan<- string, progressCh chan<- Progress) {
	var wg sync.WaitGroup
	var mutex sync.Mutex

	var totalSearched int
	var totalFound int

	ticker := time.NewTicker(1 * time.Second)

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
			return err
		}

		wg.Add(1)
		go func() {
			defer wg.Done()
			fmt.Printf("XXXX")
			// extension := strings.ToLower(filepath.Ext(path))
			// convert extension to lowercase
			if isImage(path) {
				// if extension == ".jpeg" {
				fileCh <- path
				fmt.Printf("???")
				mutex.Lock()
				fmt.Printf("!!!")
				totalFound++
				mutex.Unlock()
			}

			mutex.Lock()
			totalSearched++
			mutex.Unlock()
		}()

		return nil
	}

	fastwalk.Walk(&conf, root, walkFn)
	// filepath.Walk(root, walkFn)

	wg.Wait()
	close(fileCh)
	ticker.Stop()
}
