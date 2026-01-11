package lumina

import (
	"image"
	"image/color"
	"runtime"
	"sync"
)

// GrayscaleFilter converts an image to grayscale.
type GrayscaleFilter struct{}

// Process converts the given image to grayscale using Goroutines for performance.
func (f *GrayscaleFilter) Process(img image.Image) image.Image {
	bounds := img.Bounds()
	grayImg := image.NewGray(bounds)

	numCPU := runtime.NumCPU()
	var wg sync.WaitGroup

	// Divide the work into chunks of rows
	chunkSize := (bounds.Dy() + numCPU - 1) / numCPU

	for i := 0; i < numCPU; i++ {
		startY := bounds.Min.Y + i*chunkSize
		endY := startY + chunkSize
		if endY > bounds.Max.Y {
			endY = bounds.Max.Y
		}

		if startY >= bounds.Max.Y {
			break
		}

		wg.Add(1)
		go func(startY, endY int) {
			defer wg.Done()
			for y := startY; y < endY; y++ {
				for x := bounds.Min.X; x < bounds.Max.X; x++ {
					oldColor := img.At(x, y)
					grayColor := color.GrayModel.Convert(oldColor)
					grayImg.Set(x, y, grayColor)
				}
			}
		}(startY, endY)
	}

	wg.Wait()
	return grayImg
}
