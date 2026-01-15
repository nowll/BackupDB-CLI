package compression

import (
	"bytes"
	"compress/gzip"
	"fmt"
)

func Compress(data []byte) ([]byte, error) {
	var buf bytes.Buffer
	gzWriter := gzip.NewWriter(&buf)

	if _, err := gzWriter.Write(data); err != nil {
		return nil, fmt.Errorf("write compressed data: %w", err)
	}

	if err := gzWriter.Close(); err != nil {
		return nil, fmt.Errorf("close gzip writer: %w", err)
	}

	return buf.Bytes(), nil
}
