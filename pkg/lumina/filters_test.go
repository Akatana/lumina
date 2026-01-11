package lumina

import (
	"image"
	"image/color"
	"testing"
)

func TestGrayscaleFilter_Process(t *testing.T) {
	// Create a 10x10 RGBA image with a specific color
	img := image.NewRGBA(image.Rect(0, 0, 10, 10))
	red := color.RGBA{255, 0, 0, 255}
	for y := 0; y < 10; y++ {
		for x := 0; x < 10; x++ {
			img.Set(x, y, red)
		}
	}

	filter := &GrayscaleFilter{}
	grayImg := filter.Process(img)

	// Verify bounds
	if grayImg.Bounds() != img.Bounds() {
		t.Errorf("Expected bounds %v, got %v", img.Bounds(), grayImg.Bounds())
	}

	// Verify color conversion
	expectedGray := color.GrayModel.Convert(red).(color.Gray)
	for y := 0; y < 10; y++ {
		for x := 0; x < 10; x++ {
			gotColor := grayImg.At(x, y).(color.Gray)
			if gotColor != expectedGray {
				t.Errorf("At (%d, %d): expected %v, got %v", x, y, expectedGray, gotColor)
			}
		}
	}
}

func TestGrayscaleFilter_Process_Empty(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, 0, 0))
	filter := &GrayscaleFilter{}
	grayImg := filter.Process(img)

	if !grayImg.Bounds().Empty() {
		t.Errorf("Expected empty bounds, got %v", grayImg.Bounds())
	}
}

func TestGrayscaleFilter_Process_Large(t *testing.T) {
	// Test with image larger than NumCPU to ensure chunking logic is covered
	// Use a size that is likely to be larger than NumCPU * chunkSize
	img := image.NewRGBA(image.Rect(0, 0, 100, 100))
	filter := &GrayscaleFilter{}
	grayImg := filter.Process(img)

	if grayImg.Bounds() != img.Bounds() {
		t.Errorf("Expected bounds %v, got %v", img.Bounds(), grayImg.Bounds())
	}
}

func TestGrayscaleFilter_Process_NonZeroMin(t *testing.T) {
	// Test with image that doesn't start at (0,0)
	img := image.NewRGBA(image.Rect(5, 5, 15, 15))
	filter := &GrayscaleFilter{}
	grayImg := filter.Process(img)

	if grayImg.Bounds() != img.Bounds() {
		t.Errorf("Expected bounds %v, got %v", img.Bounds(), grayImg.Bounds())
	}
	
	// Check a pixel within bounds
	c := grayImg.At(10, 10)
	if _, ok := c.(color.Gray); !ok {
		t.Errorf("Expected color.Gray, got %T", c)
	}
}
