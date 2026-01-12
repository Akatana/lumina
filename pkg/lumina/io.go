package lumina

import (
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"path/filepath"
	"strings"

	_ "golang.org/x/image/bmp"
	_ "golang.org/x/image/webp"

	"image/gif"
	"image/jpeg"
	"image/png"

	"github.com/HugoSmits86/nativewebp"
	"golang.org/x/image/bmp"
)

// Load reads an image from the given file path.
// It supports JPEG, PNG, GIF, BMP, and WebP formats.
// Returns the decoded image, the format name (e.g., "png", "jpeg", "webp", "bmp"), and any error encountered.
func Load(path string) (image.Image, string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, "", fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	img, format, err := image.Decode(file)
	if err != nil {
		return nil, "", fmt.Errorf("failed to decode image: %w", err)
	}

	return img, format, nil
}

// Save writes an image to the given file path.
// The format is determined by the file extension. Supported extensions are:
// .jpg, .jpeg, .png, .gif, .bmp, .webp.
// Returns an error if the format is unsupported or if the file cannot be created/written.
func Save(path string, img image.Image) error {
	ext := strings.ToLower(filepath.Ext(path))

	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	switch ext {
	case ".jpg", ".jpeg":
		return jpeg.Encode(file, img, nil)
	case ".png":
		return png.Encode(file, img)
	case ".gif":
		return gif.Encode(file, img, nil)
	case ".bmp":
		return bmp.Encode(file, img)
	case ".webp":
		return nativewebp.Encode(file, img, nil)
	default:
		return fmt.Errorf("unsupported format: %s", ext)
	}
}
