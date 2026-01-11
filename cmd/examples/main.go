package main

import (
	"fmt"
	"image"
	"image/color"

	"github.com/Akatana/lumina/pkg/lumina"
)

func main() {
	// Create a simple 100x100 red image
	img := image.NewRGBA(image.Rect(0, 0, 100, 100))
	for y := 0; y < 100; y++ {
		for x := 0; x < 100; x++ {
			img.Set(x, y, color.RGBA{255, 0, 0, 255})
		}
	}

	fmt.Println("Applying Grayscale filter...")
	filter := &lumina.GrayscaleFilter{}
	grayImg := filter.Process(img)

	bounds := grayImg.Bounds()
	fmt.Printf("Processed image bounds: %v\n", bounds)

	// Check the color of a pixel
	c := grayImg.At(50, 50)
	fmt.Printf("Color at (50, 50): %v\n", c)
}
