# Lumina

A high-performance, pure-Go image processing library designed to be a CGO-free alternative for backend services and digital display systems. 🚀

### Why Lumina?

Lumina is built for developers who need fast, reliable image processing without the headaches of C dependencies.
- **Pure Go**: Zero CGO dependencies, making cross-compilation a breeze.
- **Performance Focused**: Leverages Go's concurrency primitives (Goroutines) to process images in parallel.
- **Standard Library Based**: Built on top of `image` and `image/color` for maximum compatibility.

### Installation

```bash
go get github.com/Akatana/lumina
```

### Usage

```go
import "github.com/Akatana/lumina/pkg/lumina"

func main() {
    // Load an image
    img, _, _ := lumina.Load("input.png")

    // Apply grayscale filter
    filter := &lumina.GrayscaleFilter{}
    grayImg := filter.Process(img)

    // Save the result
    lumina.Save("output.jpg", grayImg)
}
```

### Documentation

For detailed documentation, please refer to:
- [Official Go Documentation (pkg.go.dev)](https://pkg.go.dev/github.com/Akatana/lumina/pkg/lumina)
- [Project Documentation (docs/doc.md)](docs/doc.md)

### Roadmap

- [ ] Implementation of high-performance Resize and Crop algorithms.
- [ ] Support for more filters (Blur, Sharpen, etc.).
- [ ] **Dynamic asset scaling** optimized for digital signage.
- [ ] SIMD optimizations for even greater performance.
