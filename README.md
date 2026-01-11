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

// Example usage coming soon!
```

### Roadmap

- [ ] Implementation of high-performance Resize and Crop algorithms.
- [ ] Support for more filters (Blur, Sharpen, etc.).
- [ ] **Dynamic asset scaling** optimized for digital signage.
- [ ] SIMD optimizations for even greater performance.
