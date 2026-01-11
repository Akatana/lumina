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
	err := Save("test.bmp", img)
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
