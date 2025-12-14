package xutils

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"io"
)

func CompressAndBase64Encode(data []byte) (string, error) {
	var buf bytes.Buffer
	gzWritter := gzip.NewWriter(&buf)

	_, err := gzWritter.Write(data)

	if err != nil {
		return "", fmt.Errorf("Failed to compress data: %w", err)
	}

	if err := gzWritter.Close(); err != nil {
		return "", fmt.Errorf("Failed to close gzip writer: %w", err)
	}

	// ecncode to base64
	encoded := base64.StdEncoding.EncodeToString(buf.Bytes())

	return encoded, nil
}

func Base64DecodeAndDecompress(encodedData string) ([]byte, error) {
	// Decode Base64
	compressedBytes, err := base64.StdEncoding.DecodeString(encodedData)
	if err != nil {
		return nil, fmt.Errorf("invalid base64 data: %w", err)
	}

	// compress gzip
	reader, err := gzip.NewReader(bytes.NewReader(compressedBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to decompress data: %w", err)
	}
	defer reader.Close()

	// read all data
	decompressedBytes, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to read decompressed data: %w", err)
	}

	return decompressedBytes, nil
}
