package compression

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"

	"github.com/pierrec/lz4/v4"
)

type Algorithm string

const (
	Gzip Algorithm = "gzip"
	LZ4  Algorithm = "lz4"
)

type Compressor struct {
	algorithm Algorithm
}

func NewCompressor(algorithm Algorithm) *Compressor {
	return &Compressor{algorithm: algorithm}
}

func (c *Compressor) Compress(data []byte) ([]byte, error) {
	switch c.algorithm {
	case Gzip:
		return c.compressGzip(data)
	case LZ4:
		return c.compressLZ4(data)
	default:
		return nil, fmt.Errorf("unsupported algorithm: %s", c.algorithm)
	}
}

func (c *Compressor) Decompress(data []byte) ([]byte, error) {
	switch c.algorithm {
	case Gzip:
		return c.decompressGzip(data)
	case LZ4:
		return c.decompressLZ4(data)
	default:
		return nil, fmt.Errorf("unsupported algorithm: %s", c.algorithm)
	}
}

func (c *Compressor) compressGzip(data []byte) ([]byte, error) {
	var buf bytes.Buffer
	writer := gzip.NewWriter(&buf)
	if _, err := writer.Write(data); err != nil {
		return nil, err
	}
	if err := writer.Close(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (c *Compressor) decompressGzip(data []byte) ([]byte, error) {
	reader, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	defer reader.Close()
	return io.ReadAll(reader)
}

func (c *Compressor) compressLZ4(data []byte) ([]byte, error) {
	var buf bytes.Buffer
	writer := lz4.NewWriter(&buf)
	if _, err := writer.Write(data); err != nil {
		return nil, err
	}
	if err := writer.Close(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (c *Compressor) decompressLZ4(data []byte) ([]byte, error) {
	reader := lz4.NewReader(bytes.NewReader(data))
	return io.ReadAll(reader)
}
