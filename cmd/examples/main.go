package main

import (
	"fmt"
	"image"
	"image/color"
	"os"

	"github.com/Akatana/lumina/pkg/lumina"
)

func main() {
	if err := run(); err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}

func run() error {
	// Create a simple 100x100 red image
	img := image.NewRGBA(image.Rect(0, 0, 100, 100))
	for y := 0; y < 100; y++ {
		for x := 0; x < 100; x++ {
			img.Set(x, y, color.RGBA{255, 0, 0, 255})
		}
	}

	fmt.Println("Saving original image as output_orig.png...")
	err := lumina.Save("output_orig.png", img)
	if err != nil {
		return fmt.Errorf("saving image: %w", err)
	}
	defer os.Remove("output_orig.png")

	fmt.Println("Loading image back...")
	loadedImg, format, err := lumina.Load("output_orig.png")
	if err != nil {
		return fmt.Errorf("loading image: %w", err)
	}
	fmt.Printf("Loaded image format: %s\n", format)

	fmt.Println("Applying Grayscale filter...")
	filter := &lumina.GrayscaleFilter{}
	grayImg := filter.Process(loadedImg)

	fmt.Println("Resizing image to 200x200...")
	processor := &lumina.DefaultProcessor{}
	resizedImg := processor.Resize(grayImg, 200, 200)

	fmt.Println("Cropping image to 100x100...")
	croppedImg := processor.Crop(resizedImg, image.Rect(50, 50, 150, 150))

	bounds := croppedImg.Bounds()
	fmt.Printf("Processed image bounds: %v\n", bounds)

	fmt.Println("Saving processed image as output_final.jpg...")
	err = lumina.Save("output_final.jpg", croppedImg)
	if err != nil {
		return fmt.Errorf("saving processed image: %w", err)
	}
	defer os.Remove("output_final.jpg")

	// Check the color of a pixel
	c := croppedImg.At(50, 50)
	fmt.Printf("Color at (50, 50): %v\n", c)
	return nil
}
