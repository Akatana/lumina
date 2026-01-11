package lumina

import (
	"image"
	"image/color"
	"image/draw"
	"runtime"
	"sync"
)

// DefaultProcessor provides the default implementation for image operations.
// It leverages Goroutines to parallelize intensive tasks.
type DefaultProcessor struct{}

// Resize scales the image to the given dimensions using bilinear interpolation.
// It processes the image concurrently for high performance.
func (p *DefaultProcessor) Resize(img image.Image, width, height int) image.Image {
	if width <= 0 || height <= 0 {
		return image.NewRGBA(image.Rect(0, 0, 0, 0))
	}

	bounds := img.Bounds()
	srcW, srcH := bounds.Dx(), bounds.Dy()
	dst := image.NewRGBA(image.Rect(0, 0, width, height))

	numCPU := runtime.NumCPU()
	var wg sync.WaitGroup
	chunkSize := (height + numCPU - 1) / numCPU

	for i := 0; i < numCPU; i++ {
		startY := i * chunkSize
		endY := startY + chunkSize
		if endY > height {
			endY = height
		}

		if startY >= height {
			break
		}

		wg.Add(1)
		go func(startY, endY int) {
			defer wg.Done()
			for y := startY; y < endY; y++ {
				srcY := float64(y) * float64(srcH) / float64(height)
				y0 := int(srcY)
				y1 := y0 + 1
				if y1 >= srcH {
					y1 = srcH - 1
				}
				dy := srcY - float64(y0)

				for x := 0; x < width; x++ {
					srcX := float64(x) * float64(srcW) / float64(width)
					x0 := int(srcX)
					x1 := x0 + 1
					if x1 >= srcW {
						x1 = srcW - 1
					}
					dx := srcX - float64(x0)

					// Get 4 pixels for bilinear interpolation
					c00 := img.At(bounds.Min.X+x0, bounds.Min.Y+y0)
					c10 := img.At(bounds.Min.X+x1, bounds.Min.Y+y0)
					c01 := img.At(bounds.Min.X+x0, bounds.Min.Y+y1)
					c11 := img.At(bounds.Min.X+x1, bounds.Min.Y+y1)

					r00, g00, b00, a00 := c00.RGBA()
					r10, g10, b10, a10 := c10.RGBA()
					r01, g01, b01, a01 := c01.RGBA()
					r11, g11, b11, a11 := c11.RGBA()

					// Interpolate horizontally
					r0 := float64(r00)*(1-dx) + float64(r10)*dx
					g0 := float64(g00)*(1-dx) + float64(g10)*dx
					b0 := float64(b00)*(1-dx) + float64(b10)*dx
					a0 := float64(a00)*(1-dx) + float64(a10)*dx

					r1 := float64(r01)*(1-dx) + float64(r11)*dx
					g1 := float64(g01)*(1-dx) + float64(g11)*dx
					b1 := float64(b01)*(1-dx) + float64(b11)*dx
					a1 := float64(a01)*(1-dx) + float64(a11)*dx

					// Interpolate vertically
					rf := r0*(1-dy) + r1*dy
					gf := g0*(1-dy) + g1*dy
					bf := b0*(1-dy) + b1*dy
					af := a0*(1-dy) + a1*dy

					dst.SetRGBA(x, y, color.RGBA{
						R: uint8(uint32(rf) >> 8),
						G: uint8(uint32(gf) >> 8),
						B: uint8(uint32(bf) >> 8),
						A: uint8(uint32(af) >> 8),
					})
				}
			}
		}(startY, endY)
	}

	wg.Wait()
	return dst
}

// Crop extracts a rectangular sub-image from the given image.
// It uses efficient memory copying when possible.
func (p *DefaultProcessor) Crop(img image.Image, rect image.Rectangle) image.Image {
	// Intersection to ensure we don't go out of bounds
	rect = rect.Intersect(img.Bounds())
	if rect.Empty() {
		return image.NewRGBA(image.Rectangle{})
	}

	dst := image.NewRGBA(image.Rect(0, 0, rect.Dx(), rect.Dy()))
	draw.Draw(dst, dst.Bounds(), img, rect.Min, draw.Src)
	return dst
}

// ApplyFilter applies a given Filter implementation to the image.
func (p *DefaultProcessor) ApplyFilter(img image.Image, filter Filter) image.Image {
	return filter.Process(img)
}
