package lumina

import (
	"image"
	"testing"
)

// MockProcessor is a simple implementation of Processor for testing purposes.
type MockProcessor struct{}

func (p *MockProcessor) Resize(img image.Image, width, height int) image.Image {
	return image.NewRGBA(image.Rect(0, 0, width, height))
}

func (p *MockProcessor) Crop(img image.Image, rect image.Rectangle) image.Image {
	return image.NewRGBA(rect)
}

func (p *MockProcessor) ApplyFilter(img image.Image, filter Filter) image.Image {
	return filter.Process(img)
}

func TestProcessorInterface(t *testing.T) {
	var _ Processor = (*MockProcessor)(nil)

	proc := &MockProcessor{}
	img := image.NewRGBA(image.Rect(0, 0, 10, 10))

	t.Run("Resize", func(t *testing.T) {
		resized := proc.Resize(img, 20, 20)
		if resized.Bounds().Dx() != 20 || resized.Bounds().Dy() != 20 {
			t.Errorf("Expected 20x20, got %v", resized.Bounds())
		}
	})

	t.Run("Crop", func(t *testing.T) {
		rect := image.Rect(0, 0, 5, 5)
		cropped := proc.Crop(img, rect)
		if cropped.Bounds() != rect {
			t.Errorf("Expected %v, got %v", rect, cropped.Bounds())
		}
	})

	t.Run("ApplyFilter", func(t *testing.T) {
		filter := &GrayscaleFilter{}
		filtered := proc.ApplyFilter(img, filter)
		if filtered == nil {
			t.Fatal("Expected filtered image, got nil")
		}
	})
}
