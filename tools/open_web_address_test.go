package tools

import (
	"testing"
)

func TestOpenWebAddress(t *testing.T) {
	address := "https://lite.cnn.com/en"
	result := OpenWebAddress(address)
	if len(result) > 9000 || len(result) == 0 {
		t.Errorf("Expected string to be between 0 and 8000 characters, got %d", len(result))
	}
}
