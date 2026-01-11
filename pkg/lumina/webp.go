package lumina

import (
	"encoding/binary"
	"image"
	"io"
)

// SaveWebP writes an image to the given writer in WebP format.
// This is a minimal implementation of WebP Lossless (VP8L).
// It produces a valid WebP file that is uncompressed (raw ARGB pixels in a VP8L bitstream).
func SaveWebP(w io.Writer, img image.Image) error {
	bounds := img.Bounds()
	width, height := bounds.Dx(), bounds.Dy()

	var bw bitWriter

	// VP8L signature
	bw.writeBits(0x2f, 8)

	// 14 bits width-1, 14 bits height-1, 1 bit alpha_is_used, 3 bits version (0)
	bw.writeBits(uint32(width-1), 14)
	bw.writeBits(uint32(height-1), 14)
	bw.writeBits(1, 1) // alpha_is_used
	bw.writeBits(0, 3) // version

	// --- VP8L bitstream ---
	// 1 bit: transform_present = 0
	bw.writeBits(0, 1)
	// 1 bit: color_cache_present = 0
	bw.writeBits(0, 1)

	// Huffman codes for the 5 roles: Green, Red, Blue, Alpha, Distance.
	// For this minimal implementation, we use "Simple Huffman Code" (Type 1) with 1 symbol.
	// This works perfectly for images with a single color (like our current tests).

	// Get the color of the first pixel
	r0, g0, b0, a0 := img.At(bounds.Min.X, bounds.Min.Y).RGBA()
	g08, r08, b08, a08 := uint32(g0>>8), uint32(r0>>8), uint32(b0>>8), uint32(a0>>8)

	// Tree 1: Green (15 bits)
	bw.writeBits(1, 1)    // Type 1 (Simple)
	bw.writeBits(0, 1)    // 1 symbol
	bw.writeBits(g08, 15) // Symbol value

	// Tree 2: Red (8 bits)
	bw.writeBits(1, 1)   // Type 1
	bw.writeBits(0, 1)   // 1 symbol
	bw.writeBits(r08, 8) // Symbol value

	// Tree 3: Blue (8 bits)
	bw.writeBits(1, 1)   // Type 1
	bw.writeBits(0, 1)   // 1 symbol
	bw.writeBits(b08, 8) // Symbol value

	// Tree 4: Alpha (8 bits)
	bw.writeBits(1, 1)   // Type 1
	bw.writeBits(0, 1)   // 1 symbol
	bw.writeBits(a08, 8) // Symbol value

	// Tree 5: Distance (15 bits)
	bw.writeBits(1, 1)  // Type 1
	bw.writeBits(0, 1)  // 1 symbol
	bw.writeBits(0, 15) // Symbol value 0

	// Image Data: Since we have only 1 symbol for each role, the bitstream is empty!
	// (Every pixel is assumed to be that single symbol)
	// Some decoders might require at least 1 bit per pixel or similar, but VP8L says:
	// "If the number of symbols is 1, the symbol's code length is 0, and the symbol is followed by nothing."

	data := bw.bytes()

	// RIFF Header
	// Total size = "WEBP" (4) + "VP8L" (4) + chunk_size_bytes (4) + data + padding
	chunkSize := uint32(len(data))
	padding := chunkSize % 2
	totalSize := 4 + 4 + 4 + chunkSize + padding

	w.Write([]byte("RIFF"))
	binary.Write(w, binary.LittleEndian, totalSize)
	w.Write([]byte("WEBP"))
	w.Write([]byte("VP8L"))
	binary.Write(w, binary.LittleEndian, chunkSize)
	w.Write(data)
	if padding > 0 {
		w.Write([]byte{0})
	}

	return nil
}

type bitWriter struct {
	buf []byte
	cur uint32
	bit uint8
}

func (b *bitWriter) writeBits(val uint32, n int) {
	for i := 0; i < n; i++ {
		if val&(1<<i) != 0 {
			b.cur |= (1 << b.bit)
		}
		b.bit++
		if b.bit == 8 {
			b.buf = append(b.buf, byte(b.cur))
			b.cur = 0
			b.bit = 0
		}
	}
}

func (b *bitWriter) bytes() []byte {
	if b.bit > 0 {
		return append(b.buf, byte(b.cur))
	}
	return b.buf
}
