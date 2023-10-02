package validimage

import (
	"sync/atomic"

	"github.com/sourcegraph/conc/stream"
)

type Result struct {
	FilePath string
	IsValid  bool
}

type Progress struct {
	Processed      int
	ValidImages    int
	TotalFilepaths int
}

func isValidImageHandler(filePathPointer string) Result {
	isValid, _ := IsValidImage(filePathPointer)
	return Result{FilePath: filePathPointer, IsValid: isValid}
}

func ProcessImages(in []string, out chan<- Result, progress chan<- Progress) {
	processImages(in, out, progress, isValidImageHandler)
}

func processImages(
	in []string,
	out chan<- Result,
	progress chan<- Progress,
	f func(string) Result,
) {
	var total atomic.Int64
	var valid atomic.Int64

	s := stream.New().WithMaxGoroutines(10)
	for _, elem := range in {
		elem := elem
		s.Go(func() stream.Callback {
			res := f(elem)
			return func() {
				out <- res

				total.Add(1)
				if res.IsValid {
					valid.Add(1)
				}

				progress <- Progress{Processed: int(total.Load()), ValidImages: int(valid.Load()), TotalFilepaths: len(in)}
			}
		})
	}
	s.Wait()
}
