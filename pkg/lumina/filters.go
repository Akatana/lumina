package lumina

import (
	"image"
	"image/color"
	"runtime"
	"sync"
)

// GrayscaleFilter converts an image to grayscale.
// It implements the Filter interface.
type GrayscaleFilter struct{}

// Process converts the given image to grayscale using Goroutines for performance.
// It iterates over rows of the image concurrently.
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

// BlurFilter applies a simple box blur effect to the image.
type BlurFilter struct {
	Radius int
}

// Process applies the box blur effect using Goroutines.
func (f *BlurFilter) Process(img image.Image) image.Image {
	if f.Radius <= 0 {
		return img
	}

	bounds := img.Bounds()
	dst := image.NewRGBA(bounds)

	numCPU := runtime.NumCPU()
	var wg sync.WaitGroup
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
					var r, g, b, a uint32
					var count uint32

					for ky := -f.Radius; ky <= f.Radius; ky++ {
						for kx := -f.Radius; kx <= f.Radius; kx++ {
							px := x + kx
							py := y + ky
							if px >= bounds.Min.X && px < bounds.Max.X && py >= bounds.Min.Y && py < bounds.Max.Y {
								pr, pg, pb, pa := img.At(px, py).RGBA()
								r += pr
								g += pg
								b += pb
								a += pa
								count++
							}
						}
					}

					dst.SetRGBA(x, y, color.RGBA{
						R: uint8((r / count) >> 8),
						G: uint8((g / count) >> 8),
						B: uint8((b / count) >> 8),
						A: uint8((a / count) >> 8),
					})
				}
			}
		}(startY, endY)
	}

	wg.Wait()
	return dst
}

// SharpenFilter applies a sharpening effect to the image using a convolution kernel.
type SharpenFilter struct{}

// Process applies the sharpening effect using Goroutines.
func (f *SharpenFilter) Process(img image.Image) image.Image {
	bounds := img.Bounds()
	dst := image.NewRGBA(bounds)

	kernel := [3][3]float64{
		{0, -1, 0},
		{-1, 5, -1},
		{0, -1, 0},
	}

	numCPU := runtime.NumCPU()
	var wg sync.WaitGroup
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
					var fr, fg, fb, fa float64

					for ky := -1; ky <= 1; ky++ {
						for kx := -1; kx <= 1; kx++ {
							px := x + kx
							py := y + ky
							// Simple edge handling: clamp to nearest pixel
							if px < bounds.Min.X {
								px = bounds.Min.X
							} else if px >= bounds.Max.X {
								px = bounds.Max.X - 1
							}
							if py < bounds.Min.Y {
								py = bounds.Min.Y
							} else if py >= bounds.Max.Y {
								py = bounds.Max.Y - 1
							}

							pr, pg, pb, pa := img.At(px, py).RGBA()
							weight := kernel[ky+1][kx+1]
							fr += float64(pr) * weight
							fg += float64(pg) * weight
							fb += float64(pb) * weight
							fa += float64(pa) * weight
						}
					}

					dst.SetRGBA(x, y, color.RGBA{
						R: clampUint8(fr),
						G: clampUint8(fg),
						B: clampUint8(fb),
						A: clampUint8(fa),
					})
				}
			}
		}(startY, endY)
	}

	wg.Wait()
	return dst
}

// BrightnessFilter adjusts the brightness of the image.
type BrightnessFilter struct {
	Amount int // -255 to 255
}

// Process adjusts the brightness using Goroutines.
func (f *BrightnessFilter) Process(img image.Image) image.Image {
	bounds := img.Bounds()
	dst := image.NewRGBA(bounds)

	numCPU := runtime.NumCPU()
	var wg sync.WaitGroup
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
					pr, pg, pb, pa := img.At(x, y).RGBA()
					dst.SetRGBA(x, y, color.RGBA{
						R: clampUint8(float64(pr) + float64(f.Amount<<8)),
						G: clampUint8(float64(pg) + float64(f.Amount<<8)),
						B: clampUint8(float64(pb) + float64(f.Amount<<8)),
						A: uint8(pa >> 8),
					})
				}
			}
		}(startY, endY)
	}

	wg.Wait()
	return dst
}

// ContrastFilter adjusts the contrast of the image.
type ContrastFilter struct {
	Percentage float64 // -100 to 100
}

// Process adjusts the contrast using Goroutines.
func (f *ContrastFilter) Process(img image.Image) image.Image {
	bounds := img.Bounds()
	dst := image.NewRGBA(bounds)

	contrast := (100.0 + f.Percentage) / 100.0
	contrast *= contrast

	numCPU := runtime.NumCPU()
	var wg sync.WaitGroup
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
					pr, pg, pb, pa := img.At(x, y).RGBA()

					// Adjust contrast
					// (pixel - 0.5) * contrast + 0.5
					r := (((float64(pr)/65535.0)-0.5)*contrast + 0.5) * 65535.0
					g := (((float64(pg)/65535.0)-0.5)*contrast + 0.5) * 65535.0
					b := (((float64(pb)/65535.0)-0.5)*contrast + 0.5) * 65535.0

					dst.SetRGBA(x, y, color.RGBA{
						R: clampUint8(r),
						G: clampUint8(g),
						B: clampUint8(b),
						A: uint8(pa >> 8),
					})
				}
			}
		}(startY, endY)
	}

	wg.Wait()
	return dst
}

func clampUint8(v float64) uint8 {
	if v < 0 {
		return 0
	}
	if v > 65535 {
		return 255
	}
	return uint8(uint32(v) >> 8)
}
