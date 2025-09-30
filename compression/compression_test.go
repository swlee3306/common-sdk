package compression

import (
	"testing"
)

func TestGzipCompression(t *testing.T) {
	compressor := NewCompressor(Gzip)
	testData := []byte("Hello, World! This is a test message.")
	
	// Test compression
	compressed, err := compressor.Compress(testData)
	if err != nil {
		t.Fatalf("Compression failed: %v", err)
	}
	
	// Test decompression
	decompressed, err := compressor.Decompress(compressed)
	if err != nil {
		t.Fatalf("Decompression failed: %v", err)
	}
	
	// Verify data integrity
	if string(decompressed) != string(testData) {
		t.Errorf("Decompressed data doesn't match original. Got: %s, Expected: %s", 
			string(decompressed), string(testData))
	}
}

func TestLZ4Compression(t *testing.T) {
	compressor := NewCompressor(LZ4)
	testData := []byte("Hello, World! This is a test message.")
	
	// Test compression
	compressed, err := compressor.Compress(testData)
	if err != nil {
		t.Fatalf("Compression failed: %v", err)
	}
	
	// Test decompression
	decompressed, err := compressor.Decompress(compressed)
	if err != nil {
		t.Fatalf("Decompression failed: %v", err)
	}
	
	// Verify data integrity
	if string(decompressed) != string(testData) {
		t.Errorf("Decompressed data doesn't match original. Got: %s, Expected: %s", 
			string(decompressed), string(testData))
	}
}

func TestCompressionRatio(t *testing.T) {
	compressor := NewCompressor(Gzip)
	testData := []byte("This is a very long message that should compress well. " +
		"Repeating this message multiple times to test compression ratio. " +
		"This is a very long message that should compress well. " +
		"Repeating this message multiple times to test compression ratio.")
	
	compressed, err := compressor.Compress(testData)
	if err != nil {
		t.Fatalf("Compression failed: %v", err)
	}
	
	// Verify compression ratio
	originalSize := len(testData)
	compressedSize := len(compressed)
	ratio := float64(compressedSize) / float64(originalSize)
	
	if ratio >= 1.0 {
		t.Errorf("Compression ratio is not effective. Ratio: %.2f", ratio)
	}
	
	t.Logf("Compression ratio: %.2f (Original: %d, Compressed: %d)", 
		ratio, originalSize, compressedSize)
}
