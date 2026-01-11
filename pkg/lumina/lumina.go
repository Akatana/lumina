package lumina

import (
	"image"
)

// Processor defines the core operations for image processing.
type Processor interface {
	// Resize scales the image to the given dimensions.
	Resize(img image.Image, width, height int) image.Image
	// Crop extracts a sub-image from the given image.
	Crop(img image.Image, rect image.Rectangle) image.Image
	// ApplyFilter applies a filter to the image.
	ApplyFilter(img image.Image, filter Filter) image.Image
}

// Filter represents an image processing filter.
type Filter interface {
	Process(img image.Image) image.Image
}
