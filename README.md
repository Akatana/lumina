# Lumina

A high-performance, pure-Go image processing library designed to be a CGO-free alternative for backend services and digital display systems. 🚀

### Why Lumina?

Lumina is built for developers who need fast, reliable image processing without the headaches of C dependencies.
- **Pure Go**: Zero CGO dependencies, making cross-compilation a breeze.
- **Performance Focused**: Leverages Go's concurrency primitives (Goroutines) for parallel processing of Resize, Crop, and Filters.
- **Format Support**: Extensive support for PNG, JPEG, GIF, BMP, and WebP (Full decoding and encoding).
- **Standard Library Based**: Built on top of `image` and `image/color` for maximum compatibility.

### Installation

```bash
go get github.com/Akatana/lumina
```

### Usage

```go
import (
    "image"
    "github.com/Akatana/lumina/pkg/lumina"
)

func main() {
    // Load an image from a local file or a URL
    img, _, _ := lumina.Load("https://example.com/input.png")

    // Use the default processor
    processor := &lumina.DefaultProcessor{}
    
    // 1. Basic Resize
    resizedImg := processor.Resize(img, 800, 600)

    // 2. Intelligent Scaling (Perfect for Digital Signage)
    // Scale to Fill: Covers the entire area, may crop edges
    fillImg := processor.Scale(img, 1920, 1080, lumina.Fill)
    
    // Scale to Fit: Fits within area, adds letterboxing/pillarboxing
    fitImg := processor.Scale(img, 1920, 1080, lumina.Fit)

    // 3. Apply Filters
    // Grayscale
    grayImg := (&lumina.GrayscaleFilter{}).Process(fillImg)
    
    // Adjust Brightness and Contrast
    brightImg := (&lumina.BrightnessFilter{Amount: 20}).Process(grayImg)
    finalImg := (&lumina.ContrastFilter{Percentage: 10}).Process(brightImg)

    // 4. Save results
    lumina.Save("output_final.jpg", finalImg)
    lumina.Save("output_fit.png", fitImg)
}
```

### Documentation

For detailed documentation, please refer to:
- [Official Go Documentation (pkg.go.dev)](https://pkg.go.dev/github.com/Akatana/lumina/pkg/lumina)
- [Project Documentation (docs/doc.md)](docs/doc.md)

### Roadmap

- [x] Implementation of high-performance Resize and Crop algorithms.
- [x] Support for more filters (Blur, Sharpen, Brightness, Contrast).
- [x] **Dynamic asset scaling** (via `Scale` method) optimized for digital signage.
- [ ] SIMD optimizations for even greater performance.
