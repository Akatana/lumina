# Lumina Documentation

Lumina is a high-performance, pure-Go image processing library. This document provides an overview of the library's API and usage.

## Table of Contents
1. [Core Interfaces](#core-interfaces)
2. [Image I/O](#image-io)
3. [Filters](#filters)
4. [Example Usage](#example-usage)

## Core Interfaces

### Processor
The `Processor` interface defines the core operations for image processing. Lumina provides a `DefaultProcessor` implementation that uses Goroutines for high performance.

```go
type Processor interface {
    Resize(img image.Image, width, height int) image.Image
    Crop(img image.Image, rect image.Rectangle) image.Image
    ApplyFilter(img image.Image, filter Filter) image.Image
    Scale(img image.Image, targetWidth, targetHeight int, mode ScaleMode) image.Image
}
```

### ScaleMode
Defines how an image should be scaled to fit target dimensions.

- **Fit**: Scales the image to fit within the target dimensions while maintaining aspect ratio. If the aspect ratio of the target doesn't match the source, the image will be centered, and the remaining area will be filled with black (letterboxing or pillarboxing).
- **Fill**: Scales the image to completely cover the target dimensions while maintaining aspect ratio. If the aspect ratios differ, the image will be center-cropped to fit the target dimensions.
- **Stretch**: Scales the image to exactly match the target dimensions. The aspect ratio is **not** maintained, which may result in visual distortion.

#### ScaleMode Example

```go
processor := &lumina.DefaultProcessor{}

// Fit a 400x200 image into a 100x100 square
// Result: 100x100 image with a 100x50 centered content and black bars
fitImg := processor.Scale(img, 100, 100, lumina.Fit)

// Fill a 400x200 image into a 100x100 square
// Result: 100x100 image, center-cropped from the scaled 200x100 version
fillImg := processor.Scale(img, 100, 100, lumina.Fill)
```

### DefaultProcessor
`DefaultProcessor` is the standard implementation of the `Processor` interface.

- **Resize**: Uses bilinear interpolation and parallelizes row processing for speed.
- **Crop**: Efficiently extracts a sub-image using `image/draw`.
- **ApplyFilter**: Helper method to apply a `Filter` to an image.
- **Scale**: Scales an image to fit or fill a target dimension, optimized for digital signage.

### Filter
The `Filter` interface represents an image processing filter.

```go
type Filter interface {
    Process(img image.Image) image.Image
}
```

## Image I/O

Lumina provides simple functions to load and save images.

### Load
`Load(path string) (image.Image, string, error)`
Loads an image from the filesystem or a URL (http/https). Supports PNG, JPEG, GIF, BMP, and WebP.

### Save
`Save(path string, img image.Image) error`
Saves an image to the filesystem. The format is determined by the file extension:
- `.png`: Portable Network Graphics
- `.jpg`, `.jpeg`: JPEG Quality 75 (default)
- `.gif`: Graphics Interchange Format
- `.bmp`: Windows Bitmap
- `.webp`: WebP Lossless

**Note on WebP Encoding**: Lumina uses the `github.com/HugoSmits86/nativewebp` package for high-quality lossless WebP encoding. This provides full support for the WebP Lossless format without requiring CGO.

## Filters

### GrayscaleFilter
A high-performance grayscale filter that processes image rows concurrently using Goroutines.

### BlurFilter
Applies a box blur effect to the image.
- **Radius**: The radius of the blur.

```go
blurFilter := &lumina.BlurFilter{Radius: 5}
blurredImg := blurFilter.Process(img)
```

### SharpenFilter
Applies a sharpening effect using a 3x3 convolution kernel.

```go
sharpenFilter := &lumina.SharpenFilter{}
sharpenedImg := sharpenFilter.Process(img)
```

### BrightnessFilter
Adjusts the brightness of the image.
- **Amount**: The amount to adjust brightness (-255 to 255).

```go
// Increase brightness
brighterFilter := &lumina.BrightnessFilter{Amount: 50}
brighterImg := brighterFilter.Process(img)

// Decrease brightness
darkerFilter := &lumina.BrightnessFilter{Amount: -50}
darkerImg := darkerFilter.Process(img)
```

### ContrastFilter
Adjusts the contrast of the image.
- **Percentage**: The percentage to adjust contrast (-100 to 100).

```go
// Increase contrast
highContrastFilter := &lumina.ContrastFilter{Percentage: 30}
highContrastImg := highContrastFilter.Process(img)

// Decrease contrast
lowContrastFilter := &lumina.ContrastFilter{Percentage: -30}
lowContrastImg := lowContrastFilter.Process(img)
```

## Example Usage

```go
import (
    "image"
    "github.com/Akatana/lumina/pkg/lumina"
)

func main() {
    // Load an image
    img, format, err := lumina.Load("input.png")
    if err != nil {
        panic(err)
    }

    // Use DefaultProcessor for Resize and Crop
    processor := &lumina.DefaultProcessor{}
    
    // Resize to 800x600
    resizedImg := processor.Resize(img, 800, 600)

    // Crop a 400x400 area from the center (example coordinates)
    croppedImg := processor.Crop(resizedImg, image.Rect(200, 100, 600, 500))

    // Apply grayscale filter
    filter := &lumina.GrayscaleFilter{}
    grayImg := filter.Process(croppedImg)

    // Save the result
    err = lumina.Save("output.jpg", grayImg)
    if err != nil {
        panic(err)
    }
}
```

For more details, visit the official Go documentation at [pkg.go.dev/github.com/Akatana/lumina/pkg/lumina](https://pkg.go.dev/github.com/Akatana/lumina/pkg/lumina).
