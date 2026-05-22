package files

import (
	"strings"
	"testing"
)

func TestDetectMetadataReportsMimeMismatch(t *testing.T) {
	t.Parallel()

	const webpHeader = "RIFF\x1a\x00\x00\x00WEBPVP8 "

	metadata, err := DetectMetadata(".jpg", strings.NewReader(webpHeader))
	if err != nil {
		t.Fatalf("DetectMetadata returned error: %v", err)
	}

	if metadata.ExtensionType != "image/jpeg" {
		t.Fatalf("ExtensionType = %q, want image/jpeg", metadata.ExtensionType)
	}
	if metadata.DetectedType != "image/webp" {
		t.Fatalf("DetectedType = %q, want image/webp", metadata.DetectedType)
	}
	if !metadata.TypeMismatch {
		t.Fatal("TypeMismatch = false, want true")
	}
}

func TestDetectMetadataDoesNotReportMatchingMime(t *testing.T) {
	t.Parallel()

	const pngHeader = "\x89PNG\r\n\x1a\n\x00\x00\x00\rIHDR"

	metadata, err := DetectMetadata(".png", strings.NewReader(pngHeader))
	if err != nil {
		t.Fatalf("DetectMetadata returned error: %v", err)
	}

	if metadata.ExtensionType != "image/png" {
		t.Fatalf("ExtensionType = %q, want image/png", metadata.ExtensionType)
	}
	if metadata.DetectedType != "image/png" {
		t.Fatalf("DetectedType = %q, want image/png", metadata.DetectedType)
	}
	if metadata.TypeMismatch {
		t.Fatal("TypeMismatch = true, want false")
	}
}
