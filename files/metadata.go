package files

import (
	"io"
	"mime"
	"net/http"
	"strings"
)

type Metadata struct {
	ExtensionType string `json:"extensionType"`
	DetectedType  string `json:"detectedType"`
	TypeMismatch  bool   `json:"typeMismatch"`
}

func DetectMetadata(extension string, r io.Reader) (*Metadata, error) {
	buffer := make([]byte, 512)
	n, err := io.ReadFull(r, buffer)
	if err != nil && err != io.ErrUnexpectedEOF && err != io.EOF {
		return nil, err
	}

	extensionType := mime.TypeByExtension(extension)
	detectedType := http.DetectContentType(buffer[:n])

	return &Metadata{
		ExtensionType: extensionType,
		DetectedType:  detectedType,
		TypeMismatch:  mimeTypesMismatch(extensionType, detectedType),
	}, nil
}

func mimeTypesMismatch(extensionType, detectedType string) bool {
	if extensionType == "" || detectedType == "" {
		return false
	}

	return cleanMimeType(extensionType) != cleanMimeType(detectedType)
}

func cleanMimeType(value string) string {
	value = strings.ToLower(value)
	value = strings.TrimSpace(value)
	if before, _, ok := strings.Cut(value, ";"); ok {
		return strings.TrimSpace(before)
	}
	return value
}
