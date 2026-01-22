package lumina

import (
	"image"
	"image/color"
	"testing"
)

func TestDefaultProcessor(t *testing.T) {
	proc := &DefaultProcessor{}

	t.Run("Resize", func(t *testing.T) {
		// Create a 2x2 source image
		src := image.NewRGBA(image.Rect(0, 0, 2, 2))
		src.Set(0, 0, color.RGBA{255, 0, 0, 255})
		src.Set(1, 0, color.RGBA{0, 255, 0, 255})
		src.Set(0, 1, color.RGBA{0, 0, 255, 255})
		src.Set(1, 1, color.RGBA{255, 255, 255, 255})

		t.Run("Upscale", func(t *testing.T) {
			dst := proc.Resize(src, 4, 4)
			if dst.Bounds().Dx() != 4 || dst.Bounds().Dy() != 4 {
				t.Errorf("Expected 4x4, got %v", dst.Bounds())
			}
		})

		t.Run("Downscale", func(t *testing.T) {
			dst := proc.Resize(src, 1, 1)
			if dst.Bounds().Dx() != 1 || dst.Bounds().Dy() != 1 {
				t.Errorf("Expected 1x1, got %v", dst.Bounds())
			}
		})

		t.Run("Invalid Dimensions", func(t *testing.T) {
			dst := proc.Resize(src, 0, 10)
			if dst.Bounds().Dx() != 0 {
				t.Errorf("Expected width 0, got %d", dst.Bounds().Dx())
			}
			dst = proc.Resize(src, 10, -1)
			if dst.Bounds().Dy() != 0 {
				t.Errorf("Expected height 0, got %d", dst.Bounds().Dy())
			}
		})

		t.Run("Non-zero bounds", func(t *testing.T) {
			src2 := image.NewRGBA(image.Rect(10, 10, 12, 12))
			dst := proc.Resize(src2, 4, 4)
			if dst.Bounds().Dx() != 4 || dst.Bounds().Dy() != 4 {
				t.Errorf("Expected 4x4, got %v", dst.Bounds())
			}
		})
	})

	t.Run("Crop", func(t *testing.T) {
		src := image.NewRGBA(image.Rect(0, 0, 10, 10))
		for y := 0; y < 10; y++ {
			for x := 0; x < 10; x++ {
				src.Set(x, y, color.RGBA{uint8(x * 25), uint8(y * 25), 0, 255})
			}
		}

		t.Run("Valid Crop", func(t *testing.T) {
			rect := image.Rect(2, 2, 5, 5)
			dst := proc.Crop(src, rect)
			if dst.Bounds().Dx() != 3 || dst.Bounds().Dy() != 3 {
				t.Errorf("Expected 3x3, got %v", dst.Bounds())
			}
			// Check a pixel to ensure it's copied correctly
			// Note: Crop result starts at (0,0) in our implementation
			expectedColor := src.At(2, 2)
			gotColor := dst.At(0, 0)
			if gotColor != expectedColor {
				t.Errorf("Expected color %v, got %v", expectedColor, gotColor)
			}
		})

		t.Run("Out of Bounds Crop", func(t *testing.T) {
			rect := image.Rect(8, 8, 12, 12)
			dst := proc.Crop(src, rect)
			if dst.Bounds().Dx() != 2 || dst.Bounds().Dy() != 2 {
				t.Errorf("Expected 2x2, got %v", dst.Bounds())
			}
		})

		t.Run("Empty Crop", func(t *testing.T) {
			rect := image.Rect(11, 11, 15, 15)
			dst := proc.Crop(src, rect)
			if !dst.Bounds().Empty() {
				t.Errorf("Expected empty bounds, got %v", dst.Bounds())
			}
		})
	})

	t.Run("ApplyFilter", func(t *testing.T) {
		src := image.NewRGBA(image.Rect(0, 0, 5, 5))
		filter := &GrayscaleFilter{}
		dst := proc.ApplyFilter(src, filter)
		if dst == nil {
			t.Fatal("Expected non-nil image")
		}
	})

	t.Run("Scale", func(t *testing.T) {
		src := image.NewRGBA(image.Rect(0, 0, 100, 50)) // 2:1 ratio

		t.Run("Stretch", func(t *testing.T) {
			dst := proc.Scale(src, 100, 100, Stretch)
			if dst.Bounds().Dx() != 100 || dst.Bounds().Dy() != 100 {
				t.Errorf("Expected 100x100, got %v", dst.Bounds())
			}
		})

		t.Run("Fit", func(t *testing.T) {
			dst := proc.Scale(src, 200, 200, Fit)
			if dst.Bounds().Dx() != 200 || dst.Bounds().Dy() != 200 {
				t.Errorf("Expected 200x200 canvas, got %v", dst.Bounds())
			}
			// The content should be 200x100 (fitting 2:1 into 1:1)
			// We can't easily check content without complex logic, but bounds are correct
		})

		t.Run("Fill", func(t *testing.T) {
			dst := proc.Scale(src, 100, 100, Fill)
			if dst.Bounds().Dx() != 100 || dst.Bounds().Dy() != 100 {
				t.Errorf("Expected 100x100, got %v", dst.Bounds())
			}
			// 100x50 filling 100x100 will scale to 200x100 then crop to 100x100
		})
	})
}
