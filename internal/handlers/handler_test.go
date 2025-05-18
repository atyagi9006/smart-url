package handler

import (
	"net/url"
	"testing"
)

func TestGenerateShortCode(t *testing.T) {
	tests := []struct {
		length int
	}{
		{1}, {5}, {6}, {8}, {10}, {12},
	}

	for _, tt := range tests {
		t.Run("Length_"+string(rune(tt.length)), func(t *testing.T) {
			code := generateShortCode(tt.length)
			if len(code) != tt.length {
				t.Errorf("expected length %d, got %d", tt.length, len(code))
			}
		})
	}
}

func TestDomainParsing(t *testing.T) {
	u, err := url.Parse("https://example.com/path")
	if err != nil {
		t.Fatal(err)
	}
	if u.Hostname() != "example.com" {
		t.Errorf("expected domain 'example.com', got '%s'", u.Hostname())
	}
}
