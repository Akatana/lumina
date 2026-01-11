package lumina

import (
	"image"
)

// Processor defines the core operations for image processing.
// Implementations of this interface should provide methods to manipulate images.
type Processor interface {
	// Resize scales the image to the given dimensions (width and height).
	Resize(img image.Image, width, height int) image.Image
	// Crop extracts a rectangular sub-image from the given image.
	Crop(img image.Image, rect image.Rectangle) image.Image
	// ApplyFilter applies a given Filter implementation to the image.
	ApplyFilter(img image.Image, filter Filter) image.Image
}

// Filter represents an image processing filter that can be applied to an image.
type Filter interface {
	// Process takes an input image and returns a new processed image.
	Process(img image.Image) image.Image
}
