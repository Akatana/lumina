package lumina

import (
	"image"
	"image/color"
	"os"
	"path/filepath"
	"testing"
)

func TestLoadSave(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "lumina_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	img := image.NewRGBA(image.Rect(0, 0, 10, 10))
	for y := 0; y < 10; y++ {
		for x := 0; x < 10; x++ {
			img.Set(x, y, color.RGBA{255, 0, 0, 255})
		}
	}

	tests := []struct {
		name     string
		filename string
		format   string
	}{
		{"PNG", "test.png", "png"},
		{"JPEG", "test.jpg", "jpeg"},
		{"GIF", "test.gif", "gif"},
		{"BMP", "test.bmp", "bmp"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := filepath.Join(tempDir, tt.filename)

			// Test Save
			err := Save(path, img)
			if err != nil {
				t.Fatalf("Save failed: %v", err)
			}

			// Test Load
			loadedImg, format, err := Load(path)
			if err != nil {
				t.Fatalf("Load failed: %v", err)
			}

			if format != tt.format {
				t.Errorf("Expected format %s, got %s", tt.format, format)
			}

			if loadedImg.Bounds() != img.Bounds() {
				t.Errorf("Expected bounds %v, got %v", img.Bounds(), loadedImg.Bounds())
			}
		})
	}
}

func TestLoad_NonExistent(t *testing.T) {
	_, _, err := Load("non_existent_file.png")
	if err == nil {
		t.Error("Expected error for non-existent file, got nil")
	}
}

func TestLoad_Invalid(t *testing.T) {
	tempFile, err := os.CreateTemp("", "invalid_image")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tempFile.Name())

	_, err = tempFile.WriteString("not an image")
	if err != nil {
		t.Fatal(err)
	}
	tempFile.Close()

	_, _, err = Load(tempFile.Name())
	if err == nil {
		t.Error("Expected error for invalid image data, got nil")
	}
}

func TestSave_Unsupported(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, 10, 10))
	err := Save("test.unsupported", img)
	if err == nil {
		t.Error("Expected error for unsupported format, got nil")
	}
}

func TestSave_CreateFail(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, 10, 10))
	// Try to create a file in a non-existent directory
	err := Save("/non/existent/path/test.png", img)
	if err == nil {
		t.Error("Expected error for invalid path, got nil")
	}
}

func TestSaveWebP_Coverage(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, 10, 10))
	for y := 0; y < 10; y++ {
		for x := 0; x < 10; x++ {
			img.Set(x, y, color.RGBA{uint8(x * 25), uint8(y * 25), 100, 255})
		}
	}

	tempFile, err := os.CreateTemp("", "test_webp_*.webp")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tempFile.Name())
	tempFile.Close()

	err = Save(tempFile.Name(), img)
	if err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	// Test Load it back
	loadedImg, format, err := Load(tempFile.Name())
	if err != nil {
		t.Fatalf("Load WebP failed: %v", err)
	}

	if format != "webp" {
		t.Errorf("Expected format webp, got %s", format)
	}

	if loadedImg.Bounds() != img.Bounds() {
		t.Errorf("Expected bounds %v, got %v", img.Bounds(), loadedImg.Bounds())
	}

	// Verify a few pixels
	for _, p := range []image.Point{{0, 0}, {5, 5}, {9, 9}} {
		c1 := img.At(p.X, p.Y)
		c2 := loadedImg.At(p.X, p.Y)
		r1, g1, b1, a1 := c1.RGBA()
		r2, g2, b2, a2 := c2.RGBA()
		if r1 != r2 || g1 != g2 || b1 != b2 || a1 != a2 {
			t.Errorf("Pixel at %v mismatch: expected %v, got %v", p, c1, c2)
		}
	}
}
