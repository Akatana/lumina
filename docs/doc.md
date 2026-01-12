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
}
```

### DefaultProcessor
`DefaultProcessor` is the standard implementation of the `Processor` interface.

- **Resize**: Uses bilinear interpolation and parallelizes row processing for speed.
- **Crop**: Efficiently extracts a sub-image using `image/draw`.
- **ApplyFilter**: Helper method to apply a `Filter` to an image.

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
Loads an image from the filesystem. Supports PNG, JPEG, GIF, BMP, and WebP.

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
