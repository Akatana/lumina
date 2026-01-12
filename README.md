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
    // Load an image
    img, _, _ := lumina.Load("input.png")

    // Use the default processor for Resize and Crop
    processor := &lumina.DefaultProcessor{}
    
    // Resize image
    resizedImg := processor.Resize(img, 800, 600)

    // Apply grayscale filter
    filter := &lumina.GrayscaleFilter{}
    grayImg := filter.Process(resizedImg)

    // Save the result
    lumina.Save("output.jpg", grayImg)
}
```

### Documentation

For detailed documentation, please refer to:
- [Official Go Documentation (pkg.go.dev)](https://pkg.go.dev/github.com/Akatana/lumina/pkg/lumina)
- [Project Documentation (docs/doc.md)](docs/doc.md)

### Roadmap

- [x] Implementation of high-performance Resize and Crop algorithms.
- [x] Support for WebP and BMP (Load/Save).
- [ ] Support for more filters (Blur, Sharpen, etc.).
- [ ] **Dynamic asset scaling** optimized for digital signage.
- [ ] SIMD optimizations for even greater performance.
