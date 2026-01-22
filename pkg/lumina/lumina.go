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
	// Scale scales an image to fit or fill a target dimension, optimized for digital signage.
	Scale(img image.Image, targetWidth, targetHeight int, mode ScaleMode) image.Image
}

// ScaleMode defines how an image should be scaled to fit target dimensions.
type ScaleMode int

const (
	// Fit scales the image to fit within the target dimensions while maintaining aspect ratio.
	// Letterboxing or pillarboxing may occur.
	Fit ScaleMode = iota
	// Fill scales the image to cover the target dimensions while maintaining aspect ratio.
	// Parts of the image may be cropped.
	Fill
	// Stretch scales the image to exactly match the target dimensions.
	// Aspect ratio is NOT maintained.
	Stretch
)

// Filter represents an image processing filter that can be applied to an image.
type Filter interface {
	// Process takes an input image and returns a new processed image.
	Process(img image.Image) image.Image
}
