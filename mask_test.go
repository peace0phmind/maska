package mask

import (
	"regexp"
	"testing"

	"github.com/go-playground/assert/v2"
)

func TestNullMask(t *testing.T) {
	mask := NewMask(nil, nil, false)

	if mask.Masked("1a") != "1a" {
		t.Errorf("Expected '1a', got '%s'", mask.Masked("1a"))
	}
}

func TestEmptyStringMask(t *testing.T) {
	mask := NewMask("", nil, false)

	if mask.Masked("1a") != "1a" {
		t.Errorf("Expected '1a', got '%s'", mask.Masked("1a"))
	}
}

func TestUndefinedMask(t *testing.T) {
	mask := NewMask(nil, nil, false)

	if mask.Masked("1a") != "1a" {
		t.Errorf("Expected '1a', got '%s'", mask.Masked("1a"))
	}
}

func TestAtAtMask(t *testing.T) {
	mask := NewMask("@ @", nil, false)

	if mask.Masked("1") != "" {
		t.Errorf("Expected '', got '%s'", mask.Masked("1"))
	}
	if mask.Masked("a") != "a" {
		t.Errorf("Expected 'a', got '%s'", mask.Masked("a"))
	}
	if mask.Masked("ab") != "a b" {
		t.Errorf("Expected 'a b', got '%s'", mask.Masked("ab"))
	}
	if mask.Masked("abc") != "a b" {
		t.Errorf("Expected 'a b', got '%s'", mask.Masked("abc"))
	}
	if mask.Masked("1abc") != "a b" {
		t.Errorf("Expected 'a b', got '%s'", mask.Masked("1abc"))
	}

	if mask.Completed("a") != false {
		t.Errorf("Expected false, got %v", mask.Completed("a"))
	}
	if mask.Completed("ab") != true {
		t.Errorf("Expected true, got %v", mask.Completed("ab"))
	}
}

func TestCustomMask(t *testing.T) {
	tokens := Tokens{
		"D": {
			Pattern:  regexp.MustCompile(`[0-9]`),
			Optional: false,
			Multiple: false,
			Repeated: false,
		},
		"d": {
			Pattern:  regexp.MustCompile(`[0-9]`),
			Optional: true,
			Multiple: false,
			Repeated: false,
		},
		".": {
			Pattern:  regexp.MustCompile(`\.`),
			Optional: false,
			Multiple: false,
			Repeated: false,
		},
	}
	mask := NewMask("dddD.D", tokens, false)

	assert.Equal(t, mask.Masked("1234.5"), "1234.5")
	t.Log(mask.Completed("1234"))
	assert.Equal(t, mask.Masked("1.1"), "1.1")
}
